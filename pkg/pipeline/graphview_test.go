package pipeline

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

func TestGraphView(t *testing.T) {

	graphData := `{"nodes":[{"id":"80599296-81ed-427f-a3e7-aefa7b5426c6","type":"custom","initialized":false,"position":{"x":94,"y":395},"data":{"type":"int-input","name":"Integer","nodeType":"input","io":{"output":{"id":"output","valueType":"number","required":true,"readOnly":false,"value":5000000,"type":"passthrough","disabled":false}},"status":"pending","id":"80599296-81ed-427f-a3e7-aefa7b5426c6"}},{"id":"c2c686e2-6a97-40a0-bd24-97385a810996","type":"custom","initialized":false,"position":{"x":479,"y":183.99999999999997},"data":{"type":"download","name":"Download","nodeType":"downloader","io":{"filename":{"id":"filename","valueType":"string","label":"Filename","required":true,"readOnly":false,"type":"passthrough","disabled":false,"value":"1.bin"},"job":{"id":"job","valueType":"downloadjob","label":"Downloader","required":true,"readOnly":false,"type":"input","disabled":false},"limit":{"id":"limit","valueType":"number","label":"Limit","required":false,"readOnly":false,"type":"input","disabled":false,"value":5000000},"output_dir":{"id":"output_dir","valueType":"directory","label":"Output directory","required":true,"readOnly":false,"type":"input","disabled":false,"value":"./"}},"status":"pending","id":"c2c686e2-6a97-40a0-bd24-97385a810996"}},{"id":"0ace3fa3-f51d-4fbc-91fc-3a5534216a2d","type":"custom","initialized":false,"position":{"x":106,"y":124},"data":{"type":"downloader","name":"Downloader","nodeType":"downloader","io":{"job":{"id":"job","valueType":"downloadjob","label":"Downloader","required":true,"readOnly":false,"type":"output","disabled":false},"url":{"id":"url","valueType":"string","label":"Url","required":true,"readOnly":false,"type":"passthrough","disabled":false,"value":"http://localhost:9999/files/random.bin"}},"status":"pending","id":"0ace3fa3-f51d-4fbc-91fc-3a5534216a2d"}}],"edges":[{"id":"9508e572-45c0-4d07-bfbb-ee8cf507f0be","type":"default","source":"0ace3fa3-f51d-4fbc-91fc-3a5534216a2d","target":"c2c686e2-6a97-40a0-bd24-97385a810996","sourceHandle":"job","targetHandle":"job","data":{},"label":"","animated":true,"sourceX":247,"sourceY":179.84375,"targetX":475,"targetY":279.546875},{"id":"1f2e8e4e-93a4-418e-b198-cecd0d49abe3","type":"default","source":"80599296-81ed-427f-a3e7-aefa7b5426c6","target":"c2c686e2-6a97-40a0-bd24-97385a810996","sourceHandle":"output","targetHandle":"limit","data":{},"label":"","animated":true,"sourceX":235,"sourceY":434.84375,"targetX":475,"targetY":319.25}],"position":[-48.515374674441546,-5.227340253366577],"zoom":0.7578582832551992,"viewport":{"x":-48.515374674441546,"y":-5.227340253366577,"zoom":0.7578582832551992}}`

	var gv GraphView
	err := json.Unmarshal([]byte(graphData), &gv)

	if err != nil {
		t.Log(err)
		t.Fail()
	}

	registry := map[string]*Node{}
	inputNode := CreateIntInputNode()
	downloadNode := CreateDownloadNode()
	downloaderNode := CreateDownloaderNode()

	registry["int-input"] = &inputNode
	registry["download"] = &downloadNode
	registry["downloader"] = &downloaderNode

	graph := gv.ToPipelineGraph(registry)

	comm := make(chan PipelineMessage, 96)

	ctx := context.Background()
	p := NewPipeline(graph, comm)
	err = p.Run(ctx)

	if err != nil {
		t.Fail()
	}

	os.Remove("./1.bin")

}
