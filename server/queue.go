package server

type taskNode struct {
	val  *Task
	next *taskNode
}

func (tn *taskNode) add(t *Task) {

}

func (tn *taskNode) rm(id int) (t *Task) {
	for p := tn; p.next != nil; p = p.next {
		if p.next.val != nil && p.next.val.Id == id {
			t, p.next = p.next.val, p.next.next
			return
		}
	}
	return
}

func String() string {
	return ""
}
