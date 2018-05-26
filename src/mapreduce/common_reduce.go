package mapreduce

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTask int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	//
	// doReduce manages one reduce task: it should read the intermediate
	// files for the task, sort the intermediate key/value pairs by key,
	// call the user-defined reduce function (reduceF) for each key, and
	// write reduceF's output to disk.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTask) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
	// Your code here (Part I).
	//

	//
	var sliceKeys []string             //键切片集合，作为排序用
	mpKvs := make(map[string][]string) //map<key, slice>

	for i := 0; i < nMap; i++ {
		strFileName := reduceName(jobName, nMap, reduceTask) //文件名

		imm, err := os.Open(strFileName)
		if err != nil {
			log.Println("open imm file error")
		}
		var kv KeyValue
		dec := json.NewDecoder(imm)
		err = dec.Decode(&kv)
		for err != nil {
			if _, ok := mpKvs[kv.Key]; !ok {
				sliceKeys = append(sliceKeys, kv.Key)
			}
			mpKvs[kv.Key] = append(mpKvs[kv.Key], kv.Value) //插入排序

			err = dec.Decode(&kv)
		}
	}
	sort.Strings(sliceKeys)

	out, err := os.Create(outFile)
	if err != nil {
		log.Println("create output file error")
		return
	}

	encFile := json.NewEncoder(out)

	for _, key := range sliceKeys {
		err := encFile.Encode(KeyValue{key, reduceF(key, sliceKeys)})
		if err != nil {
			log.Panicln("encode file error")
		}
	}
	out.Close()
}
