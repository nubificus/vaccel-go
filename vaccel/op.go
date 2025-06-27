// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#include <vaccel.h>
*/
import "C"

type OpType int32

const (
	OpNoop                OpType = C.VACCEL_OP_NOOP
	OpBlasSgemm           OpType = C.VACCEL_OP_BLAS_SGEMM
	OpImageClassify       OpType = C.VACCEL_OP_IMAGE_CLASSIFY
	OpImageDetect         OpType = C.VACCEL_OP_IMAGE_DETECT
	OpImageSegment        OpType = C.VACCEL_OP_IMAGE_SEGMENT
	OpImagePose           OpType = C.VACCEL_OP_IMAGE_POSE
	OpImageDepth          OpType = C.VACCEL_OP_IMAGE_DEPTH
	OpExec                OpType = C.VACCEL_OP_EXEC
	OpTfModelNew          OpType = C.VACCEL_OP_TF_MODEL_NEW
	OpTfModelDestroy      OpType = C.VACCEL_OP_TF_MODEL_DESTROY
	OpTfModelRegister     OpType = C.VACCEL_OP_TF_MODEL_REGISTER
	OpTfModelUnregister   OpType = C.VACCEL_OP_TF_MODEL_UNREGISTER
	OpTfSessionLoad       OpType = C.VACCEL_OP_TF_SESSION_LOAD
	OpTfSessionRun        OpType = C.VACCEL_OP_TF_SESSION_RUN
	OpTfSessionDelete     OpType = C.VACCEL_OP_TF_SESSION_DELETE
	OpMinmax              OpType = C.VACCEL_OP_MINMAX
	OpFpgaArrayCopy       OpType = C.VACCEL_OP_FPGA_ARRAYCOPY
	OpFpgaMMult           OpType = C.VACCEL_OP_FPGA_MMULT
	OpFpgaParallel        OpType = C.VACCEL_OP_FPGA_PARALLEL
	OpFpgaVectorAdd       OpType = C.VACCEL_OP_FPGA_VECTORADD
	OpExecWithResource    OpType = C.VACCEL_OP_EXEC_WITH_RESOURCE
	OpTorchJitloadForward OpType = C.VACCEL_OP_TORCH_JITLOAD_FORWARD
	OpTorchSgemm          OpType = C.VACCEL_OP_TORCH_SGEMM
	OpOpencv              OpType = C.VACCEL_OP_OPENCV
	OpTfliteSessionLoad   OpType = C.VACCEL_OP_TFLITE_SESSION_LOAD
	OpTfliteSessionRun    OpType = C.VACCEL_OP_TFLITE_SESSION_RUN
	OpTfliteSessionDelete OpType = C.VACCEL_OP_TFLITE_SESSION_DELETE
)
