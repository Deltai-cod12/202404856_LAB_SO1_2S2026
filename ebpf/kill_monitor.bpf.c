#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

char LICENSE[] SEC("license") = "GPL";

/*
 * Tipos de eventos que Go puede recibir.
 */
#define EVENT_KILL  1
#define EVENT_PIDFD 2

/*
 * Evento enviado desde eBPF hacia Go.
 *
 * Para kill():
 *   target_pid contiene el PID destino.
 *
 * Para pidfd_send_signal():
 *   target_pid sera 0 porque el syscall recibe
 *   un pidfd y no directamente un PID.
 *   pidfd contiene el descriptor utilizado.
 */
struct kill_event {
    __u32 sender_pid;
    __u32 target_pid;
    __s32 signal;
    __s32 pidfd;
    __u32 event_type;
};

/*
 * Informacion temporal de kill().
 */
struct pending_kill {
    __u32 sender_pid;
    __u32 target_pid;
    __s32 signal;
};

/*
 * Informacion temporal de pidfd_send_signal().
 */
struct pending_pidfd {
    __u32 sender_pid;
    __s32 pidfd;
    __s32 signal;
};

/*
 * Mapa temporal para correlacionar:
 *
 * sys_enter_kill
 *       ↓
 * sys_exit_kill
 */
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, __u32);
    __type(value, struct pending_kill);
} pending_kills SEC(".maps");

/*
 * Mapa temporal para correlacionar:
 *
 * sys_enter_pidfd_send_signal
 *       ↓
 * sys_exit_pidfd_send_signal
 */
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 1024);
    __type(key, __u32);
    __type(value, struct pending_pidfd);
} pending_pidfds SEC(".maps");

/*
 * Ring Buffer utilizado para enviar eventos hacia Go.
 */
struct {
    __uint(type, BPF_MAP_TYPE_RINGBUF);
    __uint(max_entries, 256 * 1024);
} events SEC(".maps");


/* =========================================================
 * kill()
 * =========================================================
 */

SEC("tracepoint/syscalls/sys_enter_kill")
int trace_kill_enter(struct trace_event_raw_sys_enter *ctx)
{
    __u64 pid_tgid;
    __u32 tid;

    struct pending_kill pending = {};

    pid_tgid = bpf_get_current_pid_tgid();

    pending.sender_pid = pid_tgid >> 32;

    tid = (__u32)pid_tgid;

    pending.target_pid = (__u32)ctx->args[0];
    pending.signal = (__s32)ctx->args[1];

    bpf_map_update_elem(
        &pending_kills,
        &tid,
        &pending,
        BPF_ANY
    );

    return 0;
}


SEC("tracepoint/syscalls/sys_exit_kill")
int trace_kill_exit(struct trace_event_raw_sys_exit *ctx)
{
    __u64 pid_tgid;
    __u32 tid;

    struct pending_kill *pending;
    struct pending_kill data = {};

    struct kill_event *event;

    pid_tgid = bpf_get_current_pid_tgid();

    tid = (__u32)pid_tgid;

    pending = bpf_map_lookup_elem(
        &pending_kills,
        &tid
    );

    if (!pending)
        return 0;

    /*
     * Copiamos la informacion a la pila antes de
     * utilizar otros helpers eBPF.
     */
    data = *pending;

    /*
     * Solo enviamos el evento cuando kill()
     * termino correctamente.
     */
    if (ctx->ret == 0) {

        event = bpf_ringbuf_reserve(
            &events,
            sizeof(*event),
            0
        );

        if (event) {

            event->sender_pid = data.sender_pid;
            event->target_pid = data.target_pid;
            event->signal = data.signal;

            event->pidfd = -1;
            event->event_type = EVENT_KILL;

            bpf_ringbuf_submit(
                event,
                0
            );
        }
    }

    bpf_map_delete_elem(
        &pending_kills,
        &tid
    );

    return 0;
}


/* =========================================================
 * pidfd_send_signal()
 * =========================================================
 */

SEC("tracepoint/syscalls/sys_enter_pidfd_send_signal")
int trace_pidfd_enter(struct trace_event_raw_sys_enter *ctx)
{
    __u64 pid_tgid;
    __u32 tid;

    struct pending_pidfd pending = {};

    pid_tgid = bpf_get_current_pid_tgid();

    pending.sender_pid = pid_tgid >> 32;

    tid = (__u32)pid_tgid;

    /*
     * Segun el formato del tracepoint:
     *
     * args[0] = pidfd
     * args[1] = sig
     * args[2] = info
     * args[3] = flags
     */
    pending.pidfd = (__s32)ctx->args[0];
    pending.signal = (__s32)ctx->args[1];

    bpf_map_update_elem(
        &pending_pidfds,
        &tid,
        &pending,
        BPF_ANY
    );

    return 0;
}


SEC("tracepoint/syscalls/sys_exit_pidfd_send_signal")
int trace_pidfd_exit(struct trace_event_raw_sys_exit *ctx)
{
    __u64 pid_tgid;
    __u32 tid;

    struct pending_pidfd *pending;
    struct pending_pidfd data = {};

    struct kill_event *event;

    pid_tgid = bpf_get_current_pid_tgid();

    tid = (__u32)pid_tgid;

    pending = bpf_map_lookup_elem(
        &pending_pidfds,
        &tid
    );

    if (!pending)
        return 0;

    /*
     * Copiamos primero los datos encontrados
     * en el mapa.
     */
    data = *pending;

    /*
     * ret == 0 significa que pidfd_send_signal()
     * termino correctamente.
     */
    if (ctx->ret == 0) {

        event = bpf_ringbuf_reserve(
            &events,
            sizeof(*event),
            0
        );

        if (event) {

            event->sender_pid = data.sender_pid;

            /*
             * pidfd_send_signal no recibe directamente
             * el PID destino.
             */
            event->target_pid = 0;

            event->signal = data.signal;
            event->pidfd = data.pidfd;
            event->event_type = EVENT_PIDFD;

            bpf_ringbuf_submit(
                event,
                0
            );
        }
    }

    bpf_map_delete_elem(
        &pending_pidfds,
        &tid
    );

    return 0;
}