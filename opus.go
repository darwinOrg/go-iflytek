package iflytek

import "encoding/binary"

func ExtractRawOpusData(data []byte) []byte {
	var rawData []byte
	for {
		// 提取前两个字节并转换为长度
		length := int(binary.BigEndian.Uint16(data[:2]))
		if length == 0 {
			break
		}

		rawData = append(rawData, data[2:length]...)
		if len(data) > (length + 2) {
			data = data[length+2:]
		} else {
			break
		}
	}

	return rawData
}
