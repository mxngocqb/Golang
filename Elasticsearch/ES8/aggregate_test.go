package elasticsearch

import (
	"encoding/json"
	"fmt"
	"testing"
)

var aggregateClient *Client

func init() {
	var err error
	aggregateClient, err = Open("http://localhost:9200")
	if err != nil {
		panic(err)
	}
	if err := aggregateClient.Ping(); err != nil {
		panic(err)
	}
	aggregateClient.DeleteIndex("testclient_termaggregate")
	aggregateClient.DeleteIndex("testclient_rangeaggregate")
	aggregateClient.DeleteIndex("testclient_cardinalityaggregate")
	aggregateClient.DeleteIndex("testclient_compositeaggregate")
	template, _ := json.Marshal(map[string]interface{}{
		"index_patterns": []string{"*"},
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"doc": map[string]interface{}{
				"dynamic_templates": []interface{}{
					map[string]interface{}{
						"string_fields": map[string]interface{}{
							"mapping": map[string]interface{}{
								"type":  "keyword",
								"index": true,
							},
							"match_mapping_type": "string",
							"match":              "*",
						},
					},
				},
			},
		},
	})
	aggregateClient.put("_template/doc", template)
}

func TestClient_TermAggregate(t *testing.T) {
	aggregateClient.InsertDocument("testclient_termaggregate", "doc", "1", map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
	}, RefreshTrue)
	aggregateClient.InsertDocument("testclient_termaggregate", "doc", "2", map[string]interface{}{
		"field1": "value1",
		"field2": "value3",
	}, RefreshTrue)
	aggregateClient.InsertDocument("testclient_termaggregate", "doc", "3", map[string]interface{}{
		"field1": "value1",
		"field2": "value4",
	}, RefreshTrue)
	result, err := aggregateClient.TermAggregate("testclient_termaggregate", "doc", nil, NewTermAggregations([]*TermAggregation{
		{Field: "field1", Size: 10},
		{Field: "field2", Size: 10},
	}))
	if err != nil {
		t.Fatalf("could not get aggregations: %s", err)
	}
	field1 := result["field1"].Buckets
	if len(field1) != 1 || field1[0].Key.(string) != "value1" || field1[0].Count != 3 {
		t.Fatalf("wrong field1 aggs: %#v", field1)
	}
	field2 := result["field2"].Buckets
	if len(field2) != 3 || field2[0].Count != 1 || field2[1].Count != 1 || field2[2].Count != 1 {
		t.Fatalf("wrong field2 aggs: %#v", field2)
	}
}

func TestClient_RangeAggregate(t *testing.T) {
	aggregateClient.InsertDocument("testclient_rangeaggregate", "doc", "1", map[string]interface{}{"field1": 10}, RefreshTrue)
	aggregateClient.InsertDocument("testclient_rangeaggregate", "doc", "2", map[string]interface{}{"field1": 100}, RefreshTrue)
	aggregateClient.InsertDocument("testclient_rangeaggregate", "doc", "3", map[string]interface{}{"field1": 1000}, RefreshTrue)
	aggregateClient.InsertDocument("testclient_rangeaggregate", "doc", "4", map[string]interface{}{"field1": 1}, RefreshTrue)
	minValue, maxValue, err := aggregateClient.RangeAggregate("testclient_rangeaggregate", "doc", nil, "field1")
	if err != nil {
		t.Fatalf("could not range aggregate: %s", err)
	}
	if minValue != 1.0 || maxValue != 1000.0 {
		t.Fatalf("wrong range, expected %f - %f, got: %f - %f", 1.0, 1000.0, minValue, maxValue)
	}
	minValue2, maxValue2, err := aggregateClient.RangeAggregate("testclient_rangeaggregate", "doc", nil, "field2")
	if err != nil {
		t.Fatalf("could not range aggregate: %s", err)
	}
	if minValue2 != 0.0 || maxValue2 != 0.0 {
		t.Fatalf("wrong range, expected %f - %f, got: %f - %f", 0.0, 0.0, minValue2, maxValue2)
	}
}

func TestClient_CardinalityAggregate(t *testing.T) {
	aggregateClient.InsertDocument("testclient_cardinalityaggregate", "doc", "1", map[string]interface{}{"field1": 10}, RefreshTrue)
	aggregateClient.InsertDocument("testclient_cardinalityaggregate", "doc", "2", map[string]interface{}{"field1": 10}, RefreshTrue)
	aggregateClient.InsertDocument("testclient_cardinalityaggregate", "doc", "3", map[string]interface{}{"field1": 100}, RefreshTrue)
	aggregateClient.InsertDocument("testclient_cardinalityaggregate", "doc", "4", map[string]interface{}{"field1": 1}, RefreshTrue)
	value, err := aggregateClient.CardinalityAggregate("testclient_cardinalityaggregate", "doc", nil, "field1")
	if err != nil {
		t.Fatalf("could not cardinality aggregate: %s", err)
	}
	if value != 3 {
		t.Fatalf("wrong cardinality, expected 3, got: %d", value)
	}
	value2, err := aggregateClient.CardinalityAggregate("testclient_cardinalityaggregate", "doc", nil, "field2")
	if err != nil {
		t.Fatalf("could not cardinality aggregate: %s", err)
	}
	if value2 != 0 {
		t.Fatalf("wrong cardinality, expected 0, got: %d", value2)
	}
}

func TestClient_CompositeAggregate(t *testing.T) {
	for i := 0; i < 100; i++ {
		aggregateClient.InsertDocument("testclient_compositeaggregate", "doc", fmt.Sprint(i+1), map[string]interface{}{"field1": i + 1}, RefreshFalse)
	}
	aggregateClient.Refresh("testclient_compositeaggregate")
	compositeSize = 60
	buckets, err := aggregateClient.CompositeAggregate("testclient_compositeaggregate", "doc", nil, "field1")
	if err != nil {
		t.Fatalf("could not composite aggregate: %s", err)
	}
	if len(buckets) != 100 {
		t.Fatalf("wrong bucket length, expected %d, got: %d", 100, len(buckets))
	}
}
