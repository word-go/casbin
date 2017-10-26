package rbac

import (
	"testing"
	"time"
	"strconv"
	"github.com/edomosystems/casbin/util"
)

func TestSessionRole(t *testing.T) {
	rm := NewSessionRoleManager(3)
	rm.AddLink("alpha", "bravo", getCurrentTime(), getInOneHour())
	rm.AddLink("alpha", "charlie", getCurrentTime(), getInOneHour())
	rm.AddLink("bravo", "delta", getCurrentTime(), getInOneHour())
	rm.AddLink("bravo", "echo", getCurrentTime(), getInOneHour())
	rm.AddLink("charlie", "echo", getCurrentTime(), getInOneHour())
	rm.AddLink("charlie", "foxtrott", getCurrentTime(), getInOneHour())

	testSessionRole(t, rm, "alpha", "bravo", getCurrentTime(), true)
	testSessionRole(t, rm, "alpha", "charlie", getCurrentTime(), true)
	testSessionRole(t, rm, "bravo", "delta", getCurrentTime(), true)
	testSessionRole(t, rm, "bravo", "echo", getCurrentTime(), true)
	testSessionRole(t, rm, "charlie", "echo", getCurrentTime(), true)
	testSessionRole(t, rm, "charlie", "foxtrott", getCurrentTime(), true)

	testSessionRole(t, rm, "alpha", "bravo", getOneHoureAgo(), false)
	testSessionRole(t, rm, "alpha", "charlie", getOneHoureAgo(), false)
	testSessionRole(t, rm, "bravo", "delta", getOneHoureAgo(), false)
	testSessionRole(t, rm, "bravo", "echo", getOneHoureAgo(), false)
	testSessionRole(t, rm, "charlie", "echo", getOneHoureAgo(), false)
	testSessionRole(t, rm, "charlie", "foxtrott", getOneHoureAgo(), false)

	testSessionRole(t, rm, "alpha", "bravo", getInOneHour(), false)
	testSessionRole(t, rm, "alpha", "charlie", getInOneHour(), false)
	testSessionRole(t, rm, "bravo", "delta", getInOneHour(), false)
	testSessionRole(t, rm, "bravo", "echo", getInOneHour(), false)
	testSessionRole(t, rm, "charlie", "echo", getInOneHour(), false)
	testSessionRole(t, rm, "charlie", "foxtrott", getInOneHour(), false)
}

func TestAddLink(t *testing.T) {
	rm := NewSessionRoleManager(3)
	rm.AddLink("alpha", "bravo")
	testSessionRole(t, rm, "alpha", "bravo", getCurrentTime(), false)

	rm.AddLink("alpha", "bravo", getCurrentTime())
	testSessionRole(t, rm, "alpha", "bravo", getCurrentTime(), false)
}

func TestHasLink(t *testing.T) {
	rm := NewSessionRoleManager(3)

	alpha := "alpha"
	bravo := "bravo"
	if rm.HasLink(alpha, bravo) {
		t.Errorf("Role manager should not have link %s < %s", alpha, bravo)
	}
	if !rm.HasLink(alpha, alpha, getCurrentTime()) {
		t.Errorf("Role manager should have link %s < %s", alpha, alpha)
	}
	if rm.HasLink(alpha, bravo, getCurrentTime()) {
		t.Errorf("Role manager should not have link %s < %s", alpha, bravo)
	}

	rm.AddLink(alpha, bravo, getCurrentTime(), getInOneHour())
	if !rm.HasLink(alpha, bravo, getCurrentTime()) {
		t.Errorf("Role manager should have link %s < %s", alpha, bravo)
	}
}

func TestDeleteLink(t *testing.T) {
	rm := NewSessionRoleManager(3)

	alpha := "alpha"
	bravo := "bravo"
	charlie := "charlie"
	rm.AddLink(alpha, bravo, getOneHoureAgo(), getInOneHour())
	rm.AddLink(alpha, charlie, getOneHoureAgo(), getInOneHour())

	rm.DeleteLink(alpha, bravo)
	if rm.HasLink(alpha, bravo, getCurrentTime()) {
		t.Errorf("Role manager should not have link %s < %s", alpha, bravo)
	}

	rm.DeleteLink(alpha, "delta")
	rm.DeleteLink(bravo, charlie)

	if !rm.HasLink(alpha, charlie, getCurrentTime()) {
		t.Errorf("Role manager should have link %s < %s", alpha, charlie)
	}
}

func TestHierarchieLevel(t *testing.T) {
	rm := NewSessionRoleManager(2)

	rm.AddLink("alpha", "bravo", getOneHoureAgo(), getInOneHour())
	rm.AddLink("bravo", "charlie", getOneHoureAgo(), getInOneHour())
	if rm.HasLink("alpha", "charlie", getCurrentTime()) {
		t.Error("Role manager should not have link alpha < charlie")
	}
}

func TestOutdatedSessions(t *testing.T) {
	rm := NewSessionRoleManager(3)

	rm.AddLink("alpha", "bravo", getOneHoureAgo(), getCurrentTime())
	rm.AddLink("bravo", "charlie", getOneHoureAgo(), getInOneHour())

	if rm.HasLink("alpha", "bravo", getInOneHour()) {
		t.Error("Role manager should not have link alpha < bravo")
	}
	if !rm.HasLink("alpha", "charlie", getOneHoureAgo()) {
		t.Error("Role manager should have link alpha < charlie")
	}
}

func TestGetRoles(t *testing.T) {
	rm := NewSessionRoleManager(3)

	if rm.GetRoles("alpha") != nil {
		t.Error("Calling GetRoles without a time should return nil.")
	}

	if rm.GetRoles("bravo", getCurrentTime()) != nil {
		t.Error("bravo should not exist")
	}

	rm.AddLink("alpha", "bravo", getOneHoureAgo(), getInOneHour())

	testPrintSessionRoles(t, rm, "alpha", getOneHoureAgo(), []string{"bravo"})
	testPrintSessionRoles(t, rm, "alpha", getCurrentTime(), []string{"bravo"})
	testPrintSessionRoles(t, rm, "alpha", getInOneHour(), []string{})

	rm.AddLink("alpha", "charlie", getOneHoureAgo(), getCurrentTime())

	testPrintSessionRoles(t, rm, "alpha", getOneHoureAgo(), []string{"bravo", "charlie"})
	testPrintSessionRoles(t, rm, "alpha", getCurrentTime(), []string{"bravo"})
	testPrintSessionRoles(t, rm, "alpha", getInOneHour(), []string{})

	rm.AddLink("alpha", "charlie", getOneHoureAgo(), getInOneHour())

	testPrintSessionRoles(t, rm, "alpha", getOneHoureAgo(), []string{"bravo", "charlie"})
	testPrintSessionRoles(t, rm, "alpha", getCurrentTime(), []string{"bravo", "charlie"})
	testPrintSessionRoles(t, rm, "alpha", getInOneHour(), []string{})
}

func TestGetUsers(t *testing.T) {
	rm := NewSessionRoleManager(3)

	rm.AddLink("bravo", "alpha", getOneHoureAgo(), getInOneHour())
	rm.AddLink("charlie", "alpha", getOneHoureAgo(), getInOneHour())
	rm.AddLink("delta", "alpha", getOneHoureAgo(), getInOneHour())

	myRes := rm.GetUsers("alpha")
	if myRes != nil {
		t.Errorf("Calling GetUsers without a time should return nil.")
	}

	myRes = rm.GetUsers("alpha", getCurrentTime())
	if !util.ArrayEquals(myRes, []string{"bravo", "charlie", "delta"}) {
		t.Errorf("Alpha should have the using roles [bravo charlie delta] but has %s", myRes)
	}

	myRes = rm.GetUsers("alpha", getOneHoureAgo())
	if !util.ArrayEquals(myRes, []string{"bravo", "charlie", "delta"}) {
		t.Errorf("Alpha should have the using roles [bravo charlie delta] but has %s", myRes)
	}

	myRes = rm.GetUsers("alpha", getInOneHour())
	if !util.ArrayEquals(myRes, []string{}) {
		t.Errorf("Alpha should not have any using roles but has %s", myRes)
	}

	myRes = rm.GetUsers("bravo", getCurrentTime())
	if !util.ArrayEquals(myRes, []string{}) {
		t.Error("bravo should have no using roles")
	}
}

func testSessionRole(t *testing.T, rm RoleManager, name1 string, name2 string, requestTime string, res bool) {
	t.Helper()
	myRes := rm.HasLink(name1, name2, requestTime)

	if myRes != res {
		t.Errorf("%s < %s at time %s: %t, supposed to be %t", name1, name2, requestTime, !res, res)
	}
}

func testPrintSessionRoles(t *testing.T, rm RoleManager, name1 string, requestTime string, res []string) {
	t.Helper()
	myRes := rm.GetRoles(name1, requestTime)

	if !util.ArrayEquals(myRes, res) {
		t.Errorf("%s should have the roles %s at time %s, but has: %s", name1, res, requestTime, myRes)
	}
}

func getCurrentTime() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func getOneHoureAgo() string {
	return strconv.FormatInt(time.Now().UnixNano()-60*60*100000000000, 10)
}

func getInOneHour() string {
	return strconv.FormatInt(time.Now().Add(time.Hour).UnixNano(), 10)
}
