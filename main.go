package main

/*
#cgo LDFLAGS: -lmraa
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <mraa/i2c.h>

mraa_result_t
i2c_get(int bus, uint8_t device_address, uint8_t register_address, uint8_t* data, int length)
{
    mraa_result_t status = MRAA_SUCCESS;
    mraa_i2c_context i2c = mraa_i2c_init(bus);
    if (i2c == NULL) {
        return MRAA_ERROR_NO_RESOURCES;
    }
    status = mraa_i2c_address(i2c, device_address);
    if (status != MRAA_SUCCESS) {
        goto i2c_get_exit;
    }
    status = mraa_i2c_write_byte(i2c, register_address);
    if (status != MRAA_SUCCESS) {
        goto i2c_get_exit;
    }
    status = mraa_i2c_read(i2c, data, length) == length ? MRAA_SUCCESS : MRAA_ERROR_UNSPECIFIED;
    if (status != MRAA_SUCCESS) {
        goto i2c_get_exit;
    }
i2c_get_exit:
    mraa_i2c_stop(i2c);
    return status;
}
*/
import "C"
import (
	"unsafe"
	"fmt"
	"time"
	"bytes"
	"encoding/binary"
	"log"

	"github.com/HackerLoop/rotonde-client.go"
	"github.com/HackerLoop/rotonde/shared"
)


func main() {
	client := client.NewClient("ws://127.0.0.1:4224")

	eventDef := &rotonde.Definition{"PX4FLOW", "event", false, rotonde.FieldDefinitions{}}
	eventDef.PushField("frame_count", "number", "")
	eventDef.PushField("pixel_flow_x_integral", "number", "")
	eventDef.PushField("pixel_flow_y_integral", "number", "")
	eventDef.PushField("gyro_x_rate_integral", "number", "")
	eventDef.PushField("gyro_y_rate_integral", "number", "")
	eventDef.PushField("gyro_z_rate_integral", "number", "")
	eventDef.PushField("integration_timespan", "number", "")
	eventDef.PushField("sonar_timestamp", "number", "")
	eventDef.PushField("ground_distance", "number", "")
	eventDef.PushField("gyro_temperature", "number", "")
	eventDef.PushField("quality", "number", "")
	client.AddLocalDefinition(eventDef)

	var data = make([]uint8, 26)
	for {
		if C.i2c_get(1, 0x42, 0x16, (*C.uint8_t)(unsafe.Pointer(&data[0])), C.int(len(data))) != C.MRAA_SUCCESS {
			time.Sleep(1 * time.Second)
			return
		}
		fmt.Println(len(data))

		for i := 0; i < len(data); i+=1 {
			fmt.Printf("%02x", data[i])
		}
		fmt.Println()

		buffer := bytes.NewBuffer(data)
		var frame_count uint16
		var pixel_flow_x_integral int16
		var pixel_flow_y_integral int16
		var gyro_x_rate_integral int16
		var gyro_y_rate_integral int16
		var gyro_z_rate_integral int16
		var integration_timespan uint32
		var sonar_timestamp uint32
		var ground_distance int16
		var gyro_temperature int16
		var quality uint8

		if err := binary.Read(buffer, binary.LittleEndian, &frame_count); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &pixel_flow_x_integral); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &pixel_flow_y_integral); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &gyro_x_rate_integral); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &gyro_y_rate_integral); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &gyro_z_rate_integral); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &integration_timespan); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &sonar_timestamp); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &ground_distance); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &gyro_temperature); err != nil {
			log.Fatal(err)
		}
		if err := binary.Read(buffer, binary.LittleEndian, &quality); err != nil {
			log.Fatal(err)
		}

		if quality < 100 {
			continue
		}

		fmt.Println(frame_count, pixel_flow_x_integral, pixel_flow_y_integral, gyro_x_rate_integral, gyro_y_rate_integral, gyro_z_rate_integral, integration_timespan, sonar_timestamp, ground_distance, gyro_temperature, quality)
		client.SendEvent("PX4FLOW", map[string]interface{}{
			"frame_count": frame_count,
			"pixel_flow_x_integral": pixel_flow_x_integral,
			"pixel_flow_y_integral": pixel_flow_y_integral,
			"gyro_x_rate_integral": gyro_x_rate_integral,
			"gyro_y_rate_integral": gyro_y_rate_integral,
			"gyro_z_rate_integral": gyro_z_rate_integral,
			"integration_timespan": integration_timespan,
			"sonar_timestamp": sonar_timestamp,
			"ground_distance": ground_distance,
			"gyro_temperature": gyro_temperature,
			"quality": quality,
		})

		time.Sleep(50 * time.Millisecond)
	}
}
