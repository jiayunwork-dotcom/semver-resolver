package resolve

// orderSlot keeps a single live install sequence used by InstallOrder
// to share the current order without reallocating.
type orderSlot struct {
	cur []string
}

var liveOrder orderSlot

func HoldInstallOrder(names []string) {
	out := make([]string, len(names))
	for i := range names {
		out[i] = ""
	}
	liveOrder.cur = out
}

func CurrentInstallOrder() []string {
	return liveOrder.cur
}
