/*
The library provides different types of logging such as infoLog, traceLog, warningLog, errorLog and fatalLog.
The logging can be done either on standard ouput or a specified logFile.
There is a provision of logging different log types differently based on its environment meant setup.
To specify a certain log type as file logging enable its value as 1 in env var.
*/

package loglib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"project/golangLibrary/apicontext"
	"strings"
)

const (
	//INFO level 0
	INFO = iota
	//TRACE level 1
	TRACE
	//WARNING level 2
	WARNING
	//ERROR level 3
	ERROR
	//FATAL level 4
	FATAL
)

const (
	infoLogType    = "INFO_LOG_TYPE"
	traceLogType   = "TRACE_LOG_TYPE"
	warningLogType = "WARNING_LOG_TYPE"
	errorLogType   = "ERROR_LOG_TYPE"
	fatalLogType   = "FATAL_LOG_TYPE"

	logType     = "1" // "1" states file logging else console logging
	logFilePath = "LOG_FILE_DESTINATION"
)

// FieldsMap map of key value pair to log
type FieldsMap map[string]interface{}

var (
	maxLogLevel = FATAL

	info     *log.Logger
	trace    *log.Logger
	warning  *log.Logger
	errorlog *log.Logger
	fatal    *log.Logger
)

func init() {
	var (
		fileOpenErr                                       error
		fileName                                          string
		infoLog, traceLog, warningLog, errorLog, fatalLog = os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr
	)

	fileName = os.Getenv(logFilePath)

	infoLogValue := os.Getenv(infoLogType)
	if infoLogValue == logType {
		infoLog, fileOpenErr = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
		if fileName == "" {
			log.Println("file logging set but file path not found, set env var " + logFilePath)
		} else if fileOpenErr != nil {
			log.Println("could not set log file for info logs, using console logging, err: " + fileOpenErr.Error())
		}
	}

	traceLogValue := os.Getenv(traceLogType)
	if traceLogValue == logType {
		traceLog, fileOpenErr = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
		if fileName == "" {
			log.Println("file logging set but file path not found, set env var " + logFilePath)
		} else if fileOpenErr != nil {
			log.Println("could not set log file for trace logs, using console logging, err: " + fileOpenErr.Error())
		}
	}

	warningLogValue := os.Getenv(warningLogType)
	if warningLogValue == logType {
		warningLog, fileOpenErr = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
		if fileName == "" {
			log.Println("file logging set but file path not found, set env var " + logFilePath)
		} else if fileOpenErr != nil {
			log.Println("could not set log file for warning logs, using console logging, err: " + fileOpenErr.Error())
		}
	}

	errorLogValue := os.Getenv(errorLogType)
	if errorLogValue == logType {
		errorLog, fileOpenErr = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
		if fileName == "" {
			log.Println("file logging set but file path not found, set env var " + logFilePath)
		} else if fileOpenErr != nil {
			log.Println("could not set log file for error logs, using console logging, err: " + fileOpenErr.Error())
		}
	}

	fatalLogValue := os.Getenv(fatalLogType)
	if fatalLogValue == logType {
		fatalLog, fileOpenErr = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0644)
		if fileName == "" {
			log.Println("file logging set but file path not found, set env var " + logFilePath)
		} else if fileOpenErr != nil {
			log.Println("could not set log file for fatal logs, using console logging, err: " + fileOpenErr.Error())
		}
	}

	loginit(infoLog, traceLog, warningLog, errorLog, fatalLog)
}

func loginit(infoHandle, traceHandle, warningHandle, errorHandle, fatalHandle io.Writer) {
	info = log.New(infoHandle, "INFO|", log.LUTC|log.LstdFlags|log.Lshortfile)

	trace = log.New(traceHandle, "TRACE|", log.LUTC|log.LstdFlags|log.Lshortfile)

	warning = log.New(warningHandle, "WARNING|", log.LUTC|log.LstdFlags|log.Lshortfile)

	errorlog = log.New(errorHandle, "ERROR|", log.LUTC|log.LstdFlags|log.Lshortfile)

	fatal = log.New(fatalHandle, "FATAL|", log.LUTC|log.LstdFlags|log.Lshortfile)
}

func generatePrefix(ctx apicontext.CustomContext) string {
	return strings.Join([]string{ctx.UserID, ctx.UserName}, ":")
}

func generateTrackingIDs(ctx apicontext.CustomContext) string {
	requestID := ctx.RequestID
	correlationID := ctx.CorrelationID
	var retString string
	if correlationID != "" {
		retString += "correlationId=" + correlationID
	}
	retString += ":"
	if requestID != "" {
		retString += "requestId=" + requestID
	}
	return retString
}

func doLog(cLog *log.Logger, level, callDepth int, v ...interface{}) {
	if level > maxLogLevel {
		cLog.SetOutput(os.Stderr)
	}

	cLog.Output(callDepth, fmt.Sprintln(v...))
}

//Info dedicated for logging valuable information
func infoLog(v ...interface{}) {
	doLog(info, INFO, 1, v...)
}

//Trace system gives facility to helps you isolate your system problems by monitoring selected events Ex. entry and exit
func traceLog(v ...interface{}) {
	doLog(trace, TRACE, 1, v...)
}

//Warning for critical error
func warningLog(v ...interface{}) {
	doLog(warning, WARNING, 1, v...)
}

//Error logging error
func errorLog(v ...interface{}) {
	doLog(errorlog, ERROR, 1, v...)
}

//Fatal logging error
func fatalLog(v ...interface{}) {
	doLog(fatal, FATAL, 1, v...)
	os.Exit(1)
}

//GenericInfo generates info log
func GenericInfo(ctx apicontext.CustomContext, infoMessage string, fields FieldsMap) {
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	fieldsBytes, _ := json.Marshal(fields)
	fieldsString := string(fieldsBytes)
	msg := fmt.Sprintf("|%s|%s|",
		prefix,
		trackingIDs)
	if fields != nil && len(fields) > 0 {
		infoLog(msg, infoMessage, "|", fieldsString)
	} else {
		infoLog(msg, infoMessage)
	}
}

//GenericTrace generates trace log
func GenericTrace(ctx apicontext.CustomContext, traceMessage string, fields FieldsMap) {
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	msg := fmt.Sprintf("|%s|%s|",
		prefix,
		trackingIDs)
	if fields != nil && len(fields) > 0 {
		fieldsBytes, _ := json.Marshal(fields)
		fieldsString := string(fieldsBytes)
		traceLog(msg, traceMessage, "|", fieldsString)
	} else {
		traceLog(msg, traceMessage)
	}
}

//GenericWarning generates warning log
func GenericWarning(ctx apicontext.CustomContext, warnMessage string, fields FieldsMap) {
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	msg := fmt.Sprintf("|%s|%s|",
		prefix,
		trackingIDs)
	if fields != nil && len(fields) > 0 {
		fieldsBytes, _ := json.Marshal(fields)
		fieldsString := string(fieldsBytes)
		warningLog(msg, warnMessage, "|", fieldsString)
	} else {
		warningLog(msg, warnMessage)
	}
}

//GenericError generates error log
func GenericError(ctx apicontext.CustomContext, e error, fields FieldsMap) {
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	msg := ""
	if e != nil {
		msg = fmt.Sprintf("|%s|%s|%s", prefix, trackingIDs, e.Error())
	} else {
		msg = fmt.Sprintf("|%s|%s", prefix, trackingIDs)
	}

	if fields != nil && len(fields) > 0 {
		fieldsBytes, _ := json.Marshal(fields)
		fieldsString := string(fieldsBytes)
		errorLog(msg, "|", fieldsString)
	} else {
		errorLog(msg)
	}
}

//GenericFatalLog generates fatal log and then exits with os.Exit(1)
func GenericFatalLog(ctx apicontext.CustomContext, e error, fields FieldsMap) {
	prefix := generatePrefix(ctx)
	trackingIDs := generateTrackingIDs(ctx)
	msg := ""
	if e != nil {
		msg = fmt.Sprintf("|%s|%s|%s", prefix, trackingIDs, e.Error())
	} else {
		msg = fmt.Sprintf("|%s|%s", prefix, trackingIDs)
	}

	if fields != nil && len(fields) > 0 {
		fieldsBytes, _ := json.Marshal(fields)
		fieldsString := string(fieldsBytes)
		fatalLog(msg, "|", fieldsString)
	} else {
		fatalLog(msg)
	}
}
