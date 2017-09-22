package main

import (
	"atmi"
	"fmt"
	//"runtime"
	//http "net/http"
	//_ "net/http/pprof"
	"os"
	"strconv"
)

const (
	SUCCEED = 0
	FAIL    = -1
)

var M_ret int
var M_ac *atmi.ATMICtx

func assertEqual(a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	M_ac.TpLogError("%s", message)
}

//Binary main entry
func main() {

	M_ret := SUCCEED

	for i := 0; i < 10000 && SUCCEED==M_ret; i++ {

		ac, err := atmi.NewATMICtx()

		if nil != err {
			fmt.Errorf("Failed to allocate cotnext!", err)
			os.Exit(atmi.FAIL)
		}

		buf, err := ac.NewVIEW("MYVIEW1", 0)
		//		buf, err := ac.NewUBF(1024)

		if err != nil {
			ac.TpLogError("ATMI Error %s", err.Error())
			M_ret = FAIL
			goto out
		}

		s := strconv.Itoa(i)

		//Set one field for call
		if errB := buf.BVChg("tshort1", 0, s); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			goto out
		}

		if errB := buf.BVChg("tint2", 1, 123456789); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			goto out
		}

		if errB := buf.BVChg("tchar2", 4, 'C'); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			goto out
		}

		if errB := buf.BVChg("tfloat2", 0, 0.11); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			goto out
		}

		if errB := buf.BVChg("tdouble2", 0, 110.099); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			goto out
		}

		var errB1 atmi.UBFError

		if errB1 = buf.BVChg("tdouble2", 1, 110.099); nil == errB1 {
			ac.TpLogError("MUST HAVWE ERROR tdouble occ=1 does not exists, but SUCCEED!")
			M_ret = FAIL
			goto out
		}

		if errB1.Code() != atmi.BEINVAL {
			ac.TpLogError("Expeced error code %d but got %d", atmi.BEINVAL, errB1.Code())
			M_ret = FAIL
			goto out
		}

		if errB := buf.BVChg("tstring0", 2, "HELLO ENDURO"); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			goto out
		}

		b := []byte{0, 1, 2, 3, 4, 5}

		if errB := buf.BVChg("tcarray2", 0, b); nil != errB {
			ac.TpLogError("VIEW Error: %s", errB.Error())
			M_ret = FAIL
			goto out
		}

		//Call the server
		if _, err := ac.TpCall("TEST1", buf, 0); nil != err {
			ac.TpLogError("ATMI Error: %s", err.Error())
			M_ret = FAIL
			goto out
		}

		//Test the result buffer, type should be MYVIEW2

		ttshort1, errB := buf.BVGetInt16("ttshort1", 0, 0)
		assertEqual(ttshort1, 2233, "ttshort1")
		assertEqual(errB, nil, "ttshort1 -> errB")

		ttstring1, errB := buf.BVGetString("ttstring1", 0, 0)
		assertEqual(ttstring1, 2233, "ttstring1")
		assertEqual(errB, nil, "ttstring1 -> errB")

		var itype, subtype string

		// Check the buffer type & get view name
		if _, errA := ac.TpTypes(buf.Buf, &itype, &subtype); nil != errA {
			ac.TpLogError("Failed to get return buffer type: %s", errA.Message())
		}

		assertEqual(itype, "VIEW", "itype -> return buffer type")
		assertEqual(subtype, "MYVIEW2", "subtype -> return buffer type")

		ac.TpTerm()
		ac.FreeATMICtx()

		//runtime.GC()
	}

out:
	//Close the ATMI session

	os.Exit(M_ret)
}
