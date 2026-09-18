package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
)

type KillEvent struct {
	SenderPID uint32
	TargetPID uint32
	Signal    int32
}

type EBPFMonitor struct {
	collection *ebpf.Collection
	enterLink  link.Link
	exitLink   link.Link
	reader     *ringbuf.Reader
	events     chan KillEvent
}

func newEBPFMonitor() (*EBPFMonitor, error) {

	const objectPath = "../ebpf/kill_monitor.bpf.o"

	// Cargar la especificacion del archivo eBPF compilado
	spec, err := ebpf.LoadCollectionSpec(objectPath)

	if err != nil {
		return nil, fmt.Errorf(
			"error cargando especificacion eBPF: %w",
			err,
		)
	}

	// Cargar programas y mapas eBPF en el kernel
	collection, err := ebpf.NewCollection(spec)

	if err != nil {
		return nil, fmt.Errorf(
			"error cargando programa eBPF en kernel: %w",
			err,
		)
	}

	// Obtener programa de entrada de kill()
	enterProgram := collection.Programs["trace_kill_enter"]

	if enterProgram == nil {
		collection.Close()

		return nil, fmt.Errorf(
			"no se encontro el programa trace_kill_enter",
		)
	}

	// Obtener programa de salida de kill()
	exitProgram := collection.Programs["trace_kill_exit"]

	if exitProgram == nil {
		collection.Close()

		return nil, fmt.Errorf(
			"no se encontro el programa trace_kill_exit",
		)
	}

	// Conectar programa eBPF a sys_enter_kill
	enterLink, err := link.Tracepoint(
		"syscalls",
		"sys_enter_kill",
		enterProgram,
		nil,
	)

	if err != nil {
		collection.Close()

		return nil, fmt.Errorf(
			"error conectando sys_enter_kill: %w",
			err,
		)
	}

	// Conectar programa eBPF a sys_exit_kill
	exitLink, err := link.Tracepoint(
		"syscalls",
		"sys_exit_kill",
		exitProgram,
		nil,
	)

	if err != nil {
		enterLink.Close()
		collection.Close()

		return nil, fmt.Errorf(
			"error conectando sys_exit_kill: %w",
			err,
		)
	}

	// Obtener Ring Buffer creado en el programa eBPF
	eventsMap := collection.Maps["events"]

	if eventsMap == nil {
		exitLink.Close()
		enterLink.Close()
		collection.Close()

		return nil, fmt.Errorf(
			"no se encontro el mapa events",
		)
	}

	// Crear lector del Ring Buffer
	reader, err := ringbuf.NewReader(eventsMap)

	if err != nil {
		exitLink.Close()
		enterLink.Close()
		collection.Close()

		return nil, fmt.Errorf(
			"error creando ring buffer: %w",
			err,
		)
	}

	return &EBPFMonitor{
		collection: collection,
		enterLink:  enterLink,
		exitLink:   exitLink,
		reader:     reader,
		events:     make(chan KillEvent, 100),
	}, nil
}

func (monitor *EBPFMonitor) close() {

	if monitor.reader != nil {
		monitor.reader.Close()
	}

	if monitor.exitLink != nil {
		monitor.exitLink.Close()
	}

	if monitor.enterLink != nil {
		monitor.enterLink.Close()
	}

	if monitor.collection != nil {
		monitor.collection.Close()
	}
}

func (monitor *EBPFMonitor) readEvents() {

	fmt.Println(
		"Monitor eBPF esperando eventos sys_kill...",
	)

	for {

		record, err := monitor.reader.Read()

		if err != nil {

			if errors.Is(err, ringbuf.ErrClosed) {
				return
			}

			fmt.Println(
				"Error leyendo evento eBPF:",
				err,
			)

			continue
		}

		var event KillEvent

		err = binary.Read(
			bytes.NewBuffer(record.RawSample),
			binary.LittleEndian,
			&event,
		)

		if err != nil {
			fmt.Println(
				"Error decodificando evento eBPF:",
				err,
			)

			continue
		}

		/*
			Signal 0 solamente comprueba la existencia
			o permisos sobre un proceso.

			No representa una eliminacion.
		*/
		if event.Signal == 0 {
			continue
		}

		/*
			El programa eBPF solamente envia eventos
			cuando sys_exit_kill retorna 0.

			Por lo tanto, los eventos recibidos aqui
			corresponden a llamadas kill exitosas.
		*/
		monitor.events <- event
	}
}

func configureEBPFPermissions() error {

	if os.Geteuid() != 0 {
		return fmt.Errorf(
			"el daemon debe ejecutarse con sudo para cargar eBPF",
		)
	}

	return nil
}