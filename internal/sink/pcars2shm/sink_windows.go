//go:build windows

package pcars2shm

import (
	"encoding/binary"
	"fmt"
	"math"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/vitalvas/racetelemetry/internal/model"
)

const (
	shmName = "$pcars2$"
	shmSize = 65536
)

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procCreateFileMapW  = kernel32.NewProc("CreateFileMappingW")
	procMapViewOfFile   = kernel32.NewProc("MapViewOfFile")
	procUnmapViewOfFile = kernel32.NewProc("UnmapViewOfFile")
)

type shmWriter struct {
	handle syscall.Handle
	addr   uintptr
	seqNum atomic.Uint32
}

// New creates a new pCars2 shared memory sink.
func New() (*Sink, error) {
	namePtr, err := syscall.UTF16PtrFromString(shmName)
	if err != nil {
		return nil, fmt.Errorf("pcars2shm: utf16 name: %w", err)
	}

	handle, _, err := procCreateFileMapW.Call(
		uintptr(syscall.InvalidHandle),
		0,
		syscall.PAGE_READWRITE,
		0,
		uintptr(shmSize),
		uintptr(unsafe.Pointer(namePtr)),
	)

	if handle == 0 {
		return nil, fmt.Errorf("pcars2shm: CreateFileMapping: %w", err)
	}

	addr, _, err := procMapViewOfFile.Call(
		handle,
		syscall.FILE_MAP_WRITE,
		0, 0,
		uintptr(shmSize),
	)

	if addr == 0 {
		syscall.CloseHandle(syscall.Handle(handle))

		return nil, fmt.Errorf("pcars2shm: MapViewOfFile: %w", err)
	}

	return &Sink{
		writer: shmWriter{
			handle: syscall.Handle(handle),
			addr:   addr,
		},
	}, nil
}

func (w *shmWriter) write(frame *model.TelemetryFrame) error {
	buf := unsafe.Slice((*byte)(unsafe.Pointer(w.addr)), shmSize)

	// Set sequence number to odd (writing in progress)
	seq := w.seqNum.Add(1)
	binary.LittleEndian.PutUint32(buf[0:4], seq)

	// mVersion
	binary.LittleEndian.PutUint32(buf[4:8], 2)

	off := 8

	// mGameState: INGAME_PLAYING = 2
	if val(frame.IsRaceOn) {
		buf[off] = 2
	} else {
		buf[off] = 0
	}

	off++

	// mSessionState: RACE = 5
	if val(frame.IsRaceOn) {
		buf[off] = 5
	} else {
		buf[off] = 0
	}

	off++

	// mRaceState
	if val(frame.IsRaceOn) {
		buf[off] = 2 // RACING
	} else {
		buf[off] = 0
	}

	off++

	// mViewedParticipantIndex
	off++

	// mNumParticipants
	buf[off] = 1
	off++

	// Skip participant data (name, position, etc) - set basics
	off++ // padding

	// mUnfilteredThrottle
	putF32(buf, off, val(frame.Throttle))
	off += 4

	// mUnfilteredBrake
	putF32(buf, off, val(frame.Brake))
	off += 4

	// mUnfilteredSteering
	putF32(buf, off, val(frame.Steer))
	off += 4

	// mUnfilteredClutch
	putF32(buf, off, val(frame.Clutch))
	off += 4

	// mLapsInEvent
	binary.LittleEndian.PutUint16(buf[off:off+2], uint16(val(frame.TotalLaps)))
	off += 2

	// mTrackLength
	off += 4

	// mSpeed
	putF32(buf, off, val(frame.Speed))
	off += 4

	// mRpm
	putF32(buf, off, val(frame.EngineRPM))
	off += 4

	// mMaxRPM
	putF32(buf, off, val(frame.EngineMaxRPM))
	off += 4

	// mBrake
	putF32(buf, off, val(frame.Brake))
	off += 4

	// mThrottle
	putF32(buf, off, val(frame.Throttle))
	off += 4

	// mClutch
	putF32(buf, off, val(frame.Clutch))
	off += 4

	// mSteering
	putF32(buf, off, val(frame.Steer))
	off += 4

	// mGear
	buf[off] = uint8(int8(val(frame.Gear)))
	off++

	// mNumGears
	buf[off] = uint8(val(frame.NumGears))
	off++

	// mOdometerKM
	off += 4

	// mAntiLockActive
	off++

	// mBoostActive
	off++

	// mBoostAmount
	putF32(buf, off, val(frame.Boost))
	off += 4

	// mOilTempCelsius
	putF32(buf, off, val(frame.OilTemp))
	off += 4

	// mOilPressureKPa
	putF32(buf, off, val(frame.OilPressure))
	off += 4

	// mWaterTempCelsius
	putF32(buf, off, val(frame.WaterTemp))
	off += 4

	// mFuelPressureKPa
	off += 4

	// mFuelLevel
	putF32(buf, off, val(frame.Fuel))
	off += 4

	// mFuelCapacity
	putF32(buf, off, val(frame.FuelCapacity))
	off += 4

	// Tire temps [4]
	tireTemp := valArr4(frame.TireTemp)
	for i := range 4 {
		putF32(buf, off, tireTemp[i])
		off += 4
	}

	// Suspension travel [4]
	suspTravel := valArr4(frame.SuspensionTravel)
	for i := range 4 {
		putF32(buf, off, suspTravel[i])
		off += 4
	}

	// Orientation [3]
	putF32(buf, off, val(frame.Yaw))
	off += 4
	putF32(buf, off, val(frame.Pitch))
	off += 4
	putF32(buf, off, val(frame.Roll))
	off += 4

	// Local velocity [3]
	putF32(buf, off, val(frame.VelocityX))
	off += 4
	putF32(buf, off, val(frame.VelocityY))
	off += 4
	putF32(buf, off, val(frame.VelocityZ))
	off += 4

	// Angular velocity [3]
	putF32(buf, off, val(frame.AngularVelocityX))
	off += 4
	putF32(buf, off, val(frame.AngularVelocityY))
	off += 4
	putF32(buf, off, val(frame.AngularVelocityZ))
	off += 4

	// Local acceleration [3]
	putF32(buf, off, val(frame.AccelerationX))
	off += 4
	putF32(buf, off, val(frame.AccelerationY))
	off += 4
	putF32(buf, off, val(frame.AccelerationZ))
	off += 4

	// Position [3]
	putF32(buf, off, val(frame.PositionX))
	off += 4
	putF32(buf, off, val(frame.PositionY))
	off += 4
	putF32(buf, off, val(frame.PositionZ))
	off += 4

	// Set sequence number to even (write complete)
	seq = w.seqNum.Add(1)
	binary.LittleEndian.PutUint32(buf[0:4], seq)

	return nil
}

func (w *shmWriter) close() {
	if w.addr != 0 {
		procUnmapViewOfFile.Call(w.addr)
		w.addr = 0
	}

	if w.handle != 0 {
		syscall.CloseHandle(w.handle)
		w.handle = 0
	}
}

func val[T any](p *T) T {
	if p != nil {
		return *p
	}

	var zero T

	return zero
}

func valArr4[T any](p *[4]T) [4]T {
	if p != nil {
		return *p
	}

	var zero [4]T

	return zero
}

func putF32(buf []byte, off int, v float32) {
	binary.LittleEndian.PutUint32(buf[off:off+4], math.Float32bits(v))
}
