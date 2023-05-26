package fix

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (user *FxUser) constructSecurityList(session *FxSession, securityRequestID string) (string, error) {
	var securityListBody string
	var securityListParams []string
	securityListParams = append(securityListParams, formatMessageSlice(SecurityReqID, securityRequestID, true))
	securityListParams = append(securityListParams, formatMessageSlice(SecurityListRequestType, "0", true))
	securityListBody = strings.Join(securityListParams, "|")
	securityListBody = fmt.Sprintf("%s|", securityListBody)
	header := user.constructHeader(securityListBody, SecurityListRequest, session)
	headerWithBody := fmt.Sprintf("%s%s", header, securityListBody)
	trailer := constructTrailer(headerWithBody)
	message := strings.ReplaceAll(fmt.Sprintf("%s%s", headerWithBody, trailer), "|", "\u0001")
	return message, nil
}

func (user *FxUser) constructLogin(session *FxSession) (string, error) {
	var loginBody string
	var loginParams []string
	compIdSlice := strings.Split(user.SenderCompID, ".")
	if len(compIdSlice) == 1 {
		return "", errors.New(`"senderCompID" is incorrect`)
	}
	username := compIdSlice[len(compIdSlice)-1]

	loginParams = append(loginParams, formatMessageSlice(EncryptionMethod, "encryptionDisabled", false))

	loginParams = append(loginParams, formatMessageSlice(HeartbeatInterval, "30", true))

	loginParams = append(loginParams, formatMessageSlice(ResetSequence, "resetDisabled", false))

	loginParams = append(loginParams, formatMessageSlice(Username, username, true))
	loginParams = append(loginParams, formatMessageSlice(Password, user.Password, true))
	loginBody = strings.Join(loginParams, "|")
	loginBody = fmt.Sprintf("%s|", loginBody)

	header := user.constructHeader(loginBody, Logon, session)
	headerWithBody := fmt.Sprintf("%s%s", header, loginBody)
	trailer := constructTrailer(headerWithBody)
	message := strings.ReplaceAll(fmt.Sprintf("%s%s", headerWithBody, trailer), "|", "\u0001")
	return message, nil
}

func (user *FxUser) constructNewOrderSingle(session *FxSession, orderData OrderData) (string, error) {
	var newOrderSingleBody string
	var newOrderSingleParams []string
	orderData.OrderType = "market"

	volAsString := fmt.Sprintf("%g", orderData.Volume)
	fmt.Println("volume:", volAsString)
	transactTime := time.Now().UTC().Format(YYYYMMDDhhmmss)
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(ClOrdID, uuid.New().String(), true))
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(NOSSymbol, orderData.Symbol, true))
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(Side, orderData.Direction, false))
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(NOSTransactTime, transactTime, true))
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(NOSOrderQty, volAsString, true))
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(NOSOrdType, orderData.OrderType, false))
	newOrderSingleBody = strings.Join(newOrderSingleParams, "|")
	newOrderSingleBody = fmt.Sprintf("%s|", newOrderSingleBody)
	header := user.constructHeader(newOrderSingleBody, NewOrderSingle, session)
	headerWithBody := fmt.Sprintf("%s%s", header, newOrderSingleBody)
	trailer := constructTrailer(headerWithBody)
	message := strings.ReplaceAll(fmt.Sprintf("%s%s", headerWithBody, trailer), "|", "\u0001")
	return message, nil
}

func (user *FxUser) constructHeader(bodyMessage string, messageType CtraderSessionMessageType, session *FxSession) string {
	var messageTypeStr = fmt.Sprintf("%d", messageType)
	var header string
	var headerParams []string
	messageSequenceString := strconv.Itoa(session.MessageSequenceNumber)
	messageTs := time.Now().UTC().Format(YYYYMMDDhhmmss)
	headerParams = append(headerParams, formatMessageSlice(HeaderMessageType, messageTypeStr, false))
	headerParams = append(headerParams, formatMessageSlice(HeaderSenderCompId, user.SenderCompID, true))
	headerParams = append(headerParams, formatMessageSlice(HeaderTargetCompId, user.TargetCompID, true))
	headerParams = append(headerParams, formatMessageSlice(HeaderTargetSubId, "trade", false))
	headerParams = append(headerParams, formatMessageSlice(HeaderSenderSubId, user.SenderSubID, true))
	headerParams = append(headerParams, formatMessageSlice(HeaderMessageSequenceNumber, messageSequenceString, true))
	headerParams = append(headerParams, formatMessageSlice(HeaderMessageTimestamp, messageTs, true))
	header = strings.Join(headerParams, "|")
	messageLength := strconv.Itoa(len(bodyMessage) + len(header) + 1) // +1 is to account for the missing "|"

	headerParams = append([]string{formatMessageSlice(HeaderBeginString, "begin", false), formatMessageSlice(HeaderMessageLength, messageLength, true)}, headerParams...)
	header = fmt.Sprintf("%s|%s", strings.Join(headerParams, "|"), "")
	return header
}

func constructTrailer(message string) (trailer string) {
	checksumInput := strings.ReplaceAll(message, "|", "\u0001")
	checksum := strconv.Itoa(calculateChecksum(checksumInput))
	checksum = func(checksum string) string {
		if len(checksum) == 0 {
			return "000"
		}
		if len(checksum) == 1 {
			return fmt.Sprintf("00%s", checksum)
		}
		if len(checksum) == 2 {
			return fmt.Sprintf("0%s", checksum)
		}
		return checksum
	}(checksum)
	trailer = fmt.Sprintf("10=%s\u0001", checksum)
	return trailer
}

func calculateChecksum(dataToCalculate string) int {
	byteToCalculate := []byte(dataToCalculate)
	checksum := 0
	for _, chData := range byteToCalculate {
		checksum += int(chData)
	}
	return checksum % 256
}

func formatMessageSlice(ids CtraderParamIds, value string, useValueAsValue bool) string {
	if useValueAsValue {
		return fmt.Sprintf("%d=%s", ids, value)
	}
	return fmt.Sprintf("%d=%s", ids, MessageKeyValuePairs[ids][value])
}

func GetUUID() string {
	return uuid.New().String()
}
