// SPDX-License-Identifier: Apache-2.0

package vaccel

/*
#include <vaccel.h>
*/
import "C"

type OpType int32

const (
	OpNoop              OpType = C.VACCEL_OP_NOOP
	OpBlasSgemm         OpType = C.VACCEL_OP_BLAS_SGEMM
	OpImageClassify     OpType = C.VACCEL_OP_IMAGE_CLASSIFY
	OpImageDetect       OpType = C.VACCEL_OP_IMAGE_DETECT
	OpImageSegment      OpType = C.VACCEL_OP_IMAGE_SEGMENT
	OpImagePose         OpType = C.VACCEL_OP_IMAGE_POSE
	OpImageDepth        OpType = C.VACCEL_OP_IMAGE_DEPTH
	OpExec              OpType = C.VACCEL_OP_EXEC
	OpTfModelLoad       OpType = C.VACCEL_OP_TF_MODEL_LOAD
	OpTfModelUnload     OpType = C.VACCEL_OP_TF_MODEL_UNLOAD
	OpTfModelRun        OpType = C.VACCEL_OP_TF_MODEL_RUN
	OpMinmax            OpType = C.VACCEL_OP_MINMAX
	OpFpgaArrayCopy     OpType = C.VACCEL_OP_FPGA_ARRAYCOPY
	OpFpgaMMult         OpType = C.VACCEL_OP_FPGA_MMULT
	OpFpgaParallel      OpType = C.VACCEL_OP_FPGA_PARALLEL
	OpFpgaVectorAdd     OpType = C.VACCEL_OP_FPGA_VECTORADD
	OpExecWithResource  OpType = C.VACCEL_OP_EXEC_WITH_RESOURCE
	OpTorchModelLoad    OpType = C.VACCEL_OP_TORCH_MODEL_LOAD
	OpTorchModelRun     OpType = C.VACCEL_OP_TORCH_MODEL_RUN
	OpTorchSgemm        OpType = C.VACCEL_OP_TORCH_SGEMM
	OpOpencv            OpType = C.VACCEL_OP_OPENCV
	OpTfliteModelLoad   OpType = C.VACCEL_OP_TFLITE_MODEL_LOAD
	OpTfliteModelUnload OpType = C.VACCEL_OP_TFLITE_MODEL_UNLOAD
	OpTfliteModelRun    OpType = C.VACCEL_OP_TFLITE_MODEL_RUN
)
