package ztc

import (
	"sort"
	"sync"
	"testing"
	"time"
)

// resetReorderState clears the global reorder state between tests.
func resetReorderState() {
	rules.Lock()
	defer rules.Unlock()
	rules.orders = make(map[string]map[int]orderWithState)
	rules.orderer = nil
	rules.reorderDone = nil
}

// =====================================================
// Sort Logic Tests
// =====================================================

func TestSortOrders_SingleRank(t *testing.T) {
	input := map[int]orderWithState{
		100: {order: OrderRule{Order: 3, Rank: 7}},
		200: {order: OrderRule{Order: 1, Rank: 7}},
		300: {order: OrderRule{Order: 5, Rank: 7}},
		400: {order: OrderRule{Order: 2, Rank: 7}},
		500: {order: OrderRule{Order: 4, Rank: 7}},
	}

	sorted := sortOrders(input)

	if len(sorted) != 5 {
		t.Fatalf("expected 5 sorted rules, got %d", len(sorted))
	}
	expectedIDs := []int{200, 400, 100, 500, 300}
	for i, expected := range expectedIDs {
		if sorted[i].ID != expected {
			t.Errorf("position %d: expected ID %d, got %d", i, expected, sorted[i].ID)
		}
	}
}

func TestSortOrders_MixedRanks(t *testing.T) {
	input := map[int]orderWithState{
		100: {order: OrderRule{Order: 1, Rank: 7}},
		200: {order: OrderRule{Order: 1, Rank: 1}},
		300: {order: OrderRule{Order: 2, Rank: 7}},
		400: {order: OrderRule{Order: 2, Rank: 1}},
	}

	sorted := sortOrders(input)

	if len(sorted) != 4 {
		t.Fatalf("expected 4 sorted rules, got %d", len(sorted))
	}
	// Rank 1 rules come first, then rank 7
	expectedIDs := []int{200, 400, 100, 300}
	for i, expected := range expectedIDs {
		if sorted[i].ID != expected {
			t.Errorf("position %d: expected ID %d, got %d", i, expected, sorted[i].ID)
		}
	}
}

func TestRuleIDOrderPairList_Sort(t *testing.T) {
	pairs := RuleIDOrderPairList{
		{ID: 1, Order: OrderRule{Order: 5, Rank: 7}},
		{ID: 2, Order: OrderRule{Order: 2, Rank: 7}},
		{ID: 3, Order: OrderRule{Order: 1, Rank: 7}},
		{ID: 4, Order: OrderRule{Order: 3, Rank: 7}},
	}

	sort.Sort(pairs)

	expectedIDs := []int{3, 2, 4, 1}
	for i, expected := range expectedIDs {
		if pairs[i].ID != expected {
			t.Errorf("position %d: expected ID %d, got %d", i, expected, pairs[i].ID)
		}
	}
}

func TestRuleIDOrderPairList_SameOrder_TiebreakByID(t *testing.T) {
	pairs := RuleIDOrderPairList{
		{ID: 300, Order: OrderRule{Order: 1, Rank: 7}},
		{ID: 100, Order: OrderRule{Order: 1, Rank: 7}},
		{ID: 200, Order: OrderRule{Order: 1, Rank: 7}},
	}

	sort.Sort(pairs)

	expectedIDs := []int{100, 200, 300}
	for i, expected := range expectedIDs {
		if pairs[i].ID != expected {
			t.Errorf("position %d: expected ID %d, got %d", i, expected, pairs[i].ID)
		}
	}
}

// =====================================================
// Registration & Done Logic Tests
// =====================================================

func TestMarkOrderRuleAsDone(t *testing.T) {
	resetReorderState()

	rules.Lock()
	rules.orders["test_type"] = map[int]orderWithState{
		100: {order: OrderRule{Order: 1, Rank: 7}, done: false},
	}
	rules.Unlock()

	markOrderRuleAsDone(100, "test_type")

	rules.Lock()
	defer rules.Unlock()
	if !rules.orders["test_type"][100].done {
		t.Error("expected rule 100 to be marked done")
	}
}

func TestReorderWithBeforeReorder_FirstRule_RegistersWithDoneFalse(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	reorderWithBeforeReorder(
		OrderRule{Order: 1, Rank: 7}, 100, "test_reg",
		func() (int, error) { return 5, nil },
		func(id int, order OrderRule) error { return nil },
		nil,
	)

	rules.Lock()
	state, exists := rules.orders["test_reg"][100]
	ch := rules.reorderDone["test_reg"]
	rules.Unlock()

	if !exists {
		t.Fatal("expected rule 100 to be registered")
	}
	if state.done {
		t.Error("expected rule 100 to be registered with done=false")
	}
	if ch == nil {
		t.Error("expected reorderDone channel to be created")
	}

	// Cleanup: mark done so goroutine can finish
	markOrderRuleAsDone(100, "test_reg")
	waitForReorder("test_reg")
}

// =====================================================
// Full Lifecycle Tests
// =====================================================

func TestReorder_AllRulesRegisteredBeforeTick(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	var mu sync.Mutex
	reorderedIDs := map[int]int{}

	getCount := func() (int, error) { return 10, nil }
	updateOrder := func(id int, order OrderRule) error {
		mu.Lock()
		reorderedIDs[id] = order.Order
		mu.Unlock()
		return nil
	}

	// Register 5 rules rapidly (all before first tick)
	for i := 1; i <= 5; i++ {
		reorderWithBeforeReorder(
			OrderRule{Order: i, Rank: 7}, 100+i, "test_all_before",
			getCount, updateOrder, nil,
		)
	}

	// Mark all done
	for i := 1; i <= 5; i++ {
		markOrderRuleAsDone(100+i, "test_all_before")
	}

	waitForReorder("test_all_before")

	mu.Lock()
	defer mu.Unlock()

	if len(reorderedIDs) != 5 {
		t.Fatalf("expected 5 reordered rules, got %d", len(reorderedIDs))
	}
	for i := 1; i <= 5; i++ {
		if reorderedIDs[100+i] != i {
			t.Errorf("rule %d: expected order %d, got %d", 100+i, i, reorderedIDs[100+i])
		}
	}
}

func TestReorder_LateArrivingRules_NewCycleStarted(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	var mu sync.Mutex
	reorderedIDs := map[int]int{}

	getCount := func() (int, error) { return 10, nil }
	updateOrder := func(id int, order OrderRule) error {
		mu.Lock()
		reorderedIDs[id] = order.Order
		mu.Unlock()
		return nil
	}

	// Register 2 rules and mark done
	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 101, "test_late", getCount, updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 102, "test_late", getCount, updateOrder, nil)
	markOrderRuleAsDone(101, "test_late")
	markOrderRuleAsDone(102, "test_late")

	// Wait for first reorder cycle
	waitForReorder("test_late")

	mu.Lock()
	firstCycleCount := len(reorderedIDs)
	mu.Unlock()
	if firstCycleCount != 2 {
		t.Fatalf("first cycle: expected 2 reordered rules, got %d", firstCycleCount)
	}

	// Register 3 more rules (late arrivals — after first cycle completed)
	reorderWithBeforeReorder(OrderRule{Order: 3, Rank: 7}, 103, "test_late", getCount, updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 4, Rank: 7}, 104, "test_late", getCount, updateOrder, nil)
	reorderWithBeforeReorder(OrderRule{Order: 5, Rank: 7}, 105, "test_late", getCount, updateOrder, nil)
	markOrderRuleAsDone(103, "test_late")
	markOrderRuleAsDone(104, "test_late")
	markOrderRuleAsDone(105, "test_late")

	// Wait for second reorder cycle
	waitForReorder("test_late")

	mu.Lock()
	defer mu.Unlock()

	if reorderedIDs[103] != 3 {
		t.Errorf("rule 103: expected order 3, got %d", reorderedIDs[103])
	}
	if reorderedIDs[104] != 4 {
		t.Errorf("rule 104: expected order 4, got %d", reorderedIDs[104])
	}
	if reorderedIDs[105] != 5 {
		t.Errorf("rule 105: expected order 5, got %d", reorderedIDs[105])
	}
}

func TestReorder_MultipleResourceTypes_Independent(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	var mu sync.Mutex
	forwardingOrders := map[int]int{}
	dnsOrders := map[int]int{}

	getCount := func() (int, error) { return 10, nil }
	forwardingUpdate := func(id int, order OrderRule) error {
		mu.Lock()
		forwardingOrders[id] = order.Order
		mu.Unlock()
		return nil
	}
	dnsUpdate := func(id int, order OrderRule) error {
		mu.Lock()
		dnsOrders[id] = order.Order
		mu.Unlock()
		return nil
	}

	// Register traffic forwarding rules
	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 201, "forwarding_control_rule", getCount, forwardingUpdate, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 202, "forwarding_control_rule", getCount, forwardingUpdate, nil)

	// Register traffic forwarding DNS rules
	reorderWithBeforeReorder(OrderRule{Order: 1, Rank: 7}, 301, "traffic_forwarding_dns_rule", getCount, dnsUpdate, nil)
	reorderWithBeforeReorder(OrderRule{Order: 2, Rank: 7}, 302, "traffic_forwarding_dns_rule", getCount, dnsUpdate, nil)

	// Mark all done
	markOrderRuleAsDone(201, "forwarding_control_rule")
	markOrderRuleAsDone(202, "forwarding_control_rule")
	markOrderRuleAsDone(301, "traffic_forwarding_dns_rule")
	markOrderRuleAsDone(302, "traffic_forwarding_dns_rule")

	// Wait for both
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { waitForReorder("forwarding_control_rule"); wg.Done() }()
	go func() { waitForReorder("traffic_forwarding_dns_rule"); wg.Done() }()
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	if len(forwardingOrders) != 2 {
		t.Fatalf("expected 2 forwarding reordered rules, got %d", len(forwardingOrders))
	}
	if len(dnsOrders) != 2 {
		t.Fatalf("expected 2 DNS reordered rules, got %d", len(dnsOrders))
	}
	if forwardingOrders[201] != 1 {
		t.Errorf("forwarding rule 201: expected order 1, got %d", forwardingOrders[201])
	}
	if forwardingOrders[202] != 2 {
		t.Errorf("forwarding rule 202: expected order 2, got %d", forwardingOrders[202])
	}
	if dnsOrders[301] != 1 {
		t.Errorf("DNS rule 301: expected order 1, got %d", dnsOrders[301])
	}
	if dnsOrders[302] != 2 {
		t.Errorf("DNS rule 302: expected order 2, got %d", dnsOrders[302])
	}
}

func TestReorder_ConcurrentRegistration(t *testing.T) {
	resetReorderState()
	reorderTickInterval = 100 * time.Millisecond
	defer func() { reorderTickInterval = 30 * time.Second }()

	var mu sync.Mutex
	reorderedIDs := map[int]int{}

	getCount := func() (int, error) { return 10, nil }
	updateOrder := func(id int, order OrderRule) error {
		mu.Lock()
		reorderedIDs[id] = order.Order
		mu.Unlock()
		return nil
	}

	// Simulate concurrent registration (like parallelism=5)
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			reorderWithBeforeReorder(
				OrderRule{Order: idx, Rank: 7}, 400+idx, "test_concurrent",
				getCount, updateOrder, nil,
			)
			time.Sleep(50 * time.Millisecond) // simulate API call
			markOrderRuleAsDone(400+idx, "test_concurrent")
			waitForReorder("test_concurrent")
		}(i)
	}

	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	if len(reorderedIDs) != 5 {
		t.Fatalf("expected 5 reordered rules, got %d: %v", len(reorderedIDs), reorderedIDs)
	}
	for i := 1; i <= 5; i++ {
		if reorderedIDs[400+i] != i {
			t.Errorf("rule %d: expected order %d, got %d", 400+i, i, reorderedIDs[400+i])
		}
	}
}
