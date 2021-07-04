package initID

import (
	"fmt"
	"github.com/sony/sonyflake"
	"strconv"
)

var (
	sonyFlake     *sonyflake.Sonyflake
	sonyMachineID uint16
)

func getMachineID() (uint16, error) {
	return sonyMachineID, nil
}

func Init(machineID uint16) (err error) {
	sonyMachineID = machineID
	settings := sonyflake.Settings{}
	settings.MachineID = getMachineID
	sonyFlake = sonyflake.NewSonyflake(settings)
	return
}

func NewID() (id string, err error) {
	if sonyFlake == nil {
		err = fmt.Errorf("sonyFlake未初始化")
		return
	}
	//var tempID uint64
	tempID, _ := sonyFlake.NextID()
	id = strconv.FormatUint(tempID, 10)
	return
}
