#include <stdint.h>
#include <mach/mach.h>

char *PVoteCall(char *request);
void PVoteFree(char *response);
char *PVoteExportReport(char *report);

// Whole-process physical footprint, not Go heap and not an allocation estimate.
static inline uint64_t PVoteFootprintBytes(void) {
    task_vm_info_data_t info;
    mach_msg_type_number_t count = TASK_VM_INFO_COUNT;
    kern_return_t status = task_info(mach_task_self(), TASK_VM_INFO,
                                     (task_info_t)&info, &count);
    return status == KERN_SUCCESS ? info.phys_footprint : 0;
}
