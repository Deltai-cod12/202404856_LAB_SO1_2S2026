#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched/signal.h>
#include <linux/mm.h>
#include <linux/sched/mm.h>
#include <linux/sched/mm.h>
#include <linux/sched/cputime.h>
#include <linux/timekeeping.h>
#include <linux/math64.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Angel Rodriguez");
MODULE_DESCRIPTION("Modulo Kernel - Proyecto 2 SO1");
MODULE_VERSION("1.0");

#define PROC_NAME "continfo_pr2_so1_202404856"
#define CMDLINE_SIZE 256

static struct proc_dir_entry *proc_entry;

static void seq_put_json_string(struct seq_file *m, const char *str)
{
    const unsigned char *p = (const unsigned char *)str;

    seq_putc(m, '"');

    while (*p) {
        switch (*p) {
        case '"':
            seq_puts(m, "\\\"");
            break;

        case '\\':
            seq_puts(m, "\\\\");
            break;

        case '\n':
            seq_puts(m, "\\n");
            break;

        case '\r':
            seq_puts(m, "\\r");
            break;

        case '\t':
            seq_puts(m, "\\t");
            break;

        default:
            if (*p < 0x20) {
                seq_printf(m, "\\u%04x", *p);
            } else {
                seq_putc(m, *p);
            }

            break;
        }

        p++;
    }

    seq_putc(m, '"');
}

/*
 * Funcion que escribe la informacion que se mostrara
 * al leer el archivo /proc.
 */
static int proc_show(struct seq_file *m, void *v)
{
    struct sysinfo info;
    struct task_struct *task;

    unsigned long total_ram_kb;
    unsigned long free_ram_kb;
    unsigned long used_ram_kb;

    bool first_process = true;

    si_meminfo(&info);

    total_ram_kb = (info.totalram * info.mem_unit) / 1024;
    free_ram_kb = (info.freeram * info.mem_unit) / 1024;
    used_ram_kb = total_ram_kb - free_ram_kb;

    seq_puts(m, "{\n");

    seq_puts(m, "  \"ram\": {\n");
    seq_printf(m, "    \"total_kb\": %lu,\n", total_ram_kb);
    seq_printf(m, "    \"free_kb\": %lu,\n", free_ram_kb);
    seq_printf(m, "    \"used_kb\": %lu\n", used_ram_kb);
    seq_puts(m, "  },\n");

    seq_puts(m, "  \"processes\": [\n");

    for_each_process(task) {
        struct mm_struct *mm;

        unsigned long vsz_kb = 0;
        unsigned long rss_kb = 0;
        unsigned long mem_percent_x100 = 0;

        u64 utime = 0;
        u64 stime = 0;
        u64 cpu_time_ns = 0;
        u64 elapsed_ns = 0;
        u64 cpu_percent_x100 = 0;
        u64 now_ns = 0;

        char cmdline[CMDLINE_SIZE];
        int cmdline_len = 0;

        mm = get_task_mm(task);

        if (mm == NULL) {
            continue;
        }

        /*
         * Memoria del proceso.
         */
        vsz_kb = (mm->total_vm << PAGE_SHIFT) / 1024;
        rss_kb = get_mm_rss(mm) << (PAGE_SHIFT - 10);

        if (total_ram_kb > 0) {
            mem_percent_x100 =
                (rss_kb * 10000UL) / total_ram_kb;
        }

        /*
         * Command Line.
         */
        memset(cmdline, 0, sizeof(cmdline));

        if (mm->arg_end > mm->arg_start) {
            unsigned long arg_length =
                mm->arg_end - mm->arg_start;

            if (arg_length > CMDLINE_SIZE - 1) {
                arg_length = CMDLINE_SIZE - 1;
            }

            cmdline_len = access_process_vm(
                task,
                mm->arg_start,
                cmdline,
                arg_length,
                FOLL_FORCE
            );

            if (cmdline_len > 0) {
                int i;

                for (i = 0; i < cmdline_len; i++) {
                    if (cmdline[i] == '\0') {
                        cmdline[i] = ' ';
                    }
                }

                cmdline[cmdline_len] = '\0';
            }
        }

        if (cmdline_len <= 0) {
            strcpy(cmdline, "N/A");
        }

        /*
         * CPU promedio del proceso.
         */
        task_cputime_adjusted(task, &utime, &stime);

        cpu_time_ns = utime + stime;
        now_ns = ktime_get_ns();

        if (now_ns > task->start_time) {
            elapsed_ns = now_ns - task->start_time;
        }

        if (elapsed_ns > 0) {
            cpu_percent_x100 =
                div64_u64(
                    cpu_time_ns * 10000ULL,
                    elapsed_ns
                );
        }

        /*
         * Separar objetos JSON mediante coma.
         */
        if (!first_process) {
            seq_puts(m, ",\n");
        }

        first_process = false;

        seq_puts(m, "    {\n");

        seq_printf(
            m,
            "      \"pid\": %d,\n",
            task->pid
        );

        seq_puts(m, "      \"name\": ");
        seq_put_json_string(m, task->comm);
        seq_puts(m, ",\n");

        seq_puts(m, "      \"cmd\": ");
        seq_put_json_string(m, cmdline);
        seq_puts(m, ",\n");

        seq_printf(
            m,
            "      \"vsz_kb\": %lu,\n",
            vsz_kb
        );

        seq_printf(
            m,
            "      \"rss_kb\": %lu,\n",
            rss_kb
        );

        seq_printf(
            m,
            "      \"mem_percent_x100\": %lu,\n",
            mem_percent_x100
        );

        seq_printf(
            m,
            "      \"cpu_percent_x100\": %llu\n",
            cpu_percent_x100
        );

        seq_puts(m, "    }");

        mmput(mm);
    }

    seq_puts(m, "\n  ]\n");
    seq_puts(m, "}\n");

    return 0;
}

/*
 * Funcion ejecutada cuando se abre el archivo /proc.
 */
static int proc_open(struct inode *inode, struct file *file)
{
    return single_open(file, proc_show, NULL);
}

/*
 * Operaciones permitidas sobre el archivo /proc.
 */
static const struct proc_ops proc_fops = {
    .proc_open = proc_open,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};

/*
 * Se ejecuta al cargar el modulo.
 */
static int __init modulo_init(void)
{
    proc_entry = proc_create(
        PROC_NAME,
        0444,
        NULL,
        &proc_fops
    );

    if (proc_entry == NULL) {
        printk(KERN_ERR
               "SO1_Proyecto2: Error al crear /proc/%s\n",
               PROC_NAME);

        return -ENOMEM;
    }

    printk(KERN_INFO
           "SO1_Proyecto2: Modulo cargado correctamente - 202404856\n");

    printk(KERN_INFO
           "SO1_Proyecto2: /proc/%s creado correctamente\n",
           PROC_NAME);

    return 0;
}

/*
 * Se ejecuta al descargar el modulo.
 */
static void __exit modulo_exit(void)
{
    proc_remove(proc_entry);

    printk(KERN_INFO
           "SO1_Proyecto2: /proc/%s eliminado correctamente\n",
           PROC_NAME);

    printk(KERN_INFO
           "SO1_Proyecto2: Modulo descargado correctamente - 202404856\n");
}

module_init(modulo_init);
module_exit(modulo_exit);