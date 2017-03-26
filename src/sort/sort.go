package main

import (
	"log"
	"os"
	"bufio"
	"io"
	// "fmt"
	"time"
	"sort"
	"strconv"
	"os/exec"
	"errors"
)

// 最大运行放入内存的数量
const MAX_MEM_SIZE = 100
// 每个大组包含的小分组数
const PER_GROUP_SIZE = 25
// 每次从每个小组取出的数量
const PER_GROUP_NUM = 4

// 调整参数时，确保 PER_GROUP_SIZE * PER_GROUP_NUM <= MAX_MEM_SIZE

func main() {
	startTime := time.Now().Unix()
	f, err := os.OpenFile("test_data.in", os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	cmdObj2 := exec.Command("/bin/sh", "-c", "rm test_data.out*")
	cmdObj2.Output()

	/*
	var a []byte
	a = make([]byte, 1)
	n, err := f.ReadAt(a, 0)
	log.Println(n, err, string(a), len(a))
	*/

	line := 0
	reader := bufio.NewReader(f)
	var arr []int
	arrIdx := 0
	var files []string

	log.Println("Config: ", MAX_MEM_SIZE, PER_GROUP_SIZE, PER_GROUP_NUM)
	log.Println("Start Divide file into small files")
	for {
		// str, err := reader.ReadString('\n')
		str, _, err := reader.ReadLine()

		if err == io.EOF || len(arr) >= MAX_MEM_SIZE {
			// 开始排序
			if len(arr) > 0 {
				sortFileLog(arr)
				f_name := "test_data.out." + strconv.Itoa(arrIdx)
				files = append(files, f_name)
				f_tmp, _ := os.Create(f_name)
				for _, t_d := range(arr) {
					t_d := strconv.Itoa(t_d)
					f_tmp.WriteString(t_d + "\n")
				}
				arrIdx++
				arr = arr[0:0]
			}
			if err == io.EOF {
				// log.Println(line, err)
				break
			}
		}

		// fmt.Printf("%d: %s\n", line, str)
		if len(str) != 0 {
			tmp_int, _ := strconv.Atoi(string(str))
			arr = append(arr, tmp_int)
			
		}
		line++
		// time.Sleep(time.Second)
	}

	log.Println("Done Divide file into small files, ", len(files))

	file_name_prefix := 0
	for {
		log.Println("Start File Merge sort round: ", file_name_prefix)

		st := 0
		tmp_len := len(files)

		var new_files []string
		// var arr_files []([]string)
		for {
			tmp_start := st*PER_GROUP_SIZE
			if tmp_start >= tmp_len {
				break
			}
			tmp_end := (st+1)*PER_GROUP_SIZE
			if tmp_end > tmp_len {
				tmp_end = tmp_len
			}

			tmp_files := files[tmp_start:tmp_end]
			if len(tmp_files) == 0 {
				break
			}

			var t T
			t.fileReaders = make(fileReaders1)
			t.fileNumBuffers = make(fileNumBuffers1)

			t.openFiles(tmp_files)

			ff_name := "test_data.out." + strconv.Itoa(file_name_prefix) + "." + strconv.Itoa(st)
			// log.Println(ff_name)
			ff_tmp, _ := os.Create(ff_name)

			new_files = append(new_files, ff_name)

			for {
				min_num, err := t.fetchNumFromFiles()
				// log.Println(min_num, "aa", err)
				if err != nil {
					break;
				}
				
				ff_tmp.Seek(0, os.SEEK_END)
				ff_tmp.WriteString(strconv.Itoa(min_num) + "\n")
			}

			// arr_files = append(arr_files, tmp_files)
			st++
		}
		for _, to_be_del := range(files) {
			os.Remove(to_be_del)
		}

		log.Println("Done File Merge sort round: ", file_name_prefix, ", files: ", len(new_files))

		if len(new_files) <= 1 {
			log.Println("Sort finished", new_files)
			break;
		}

		files = new_files

		file_name_prefix++
	}

	endTime := time.Now().Unix()

	log.Println("Total Cost: ", (endTime - startTime) , "s")
}

type T struct {
	fileReaders fileReaders1
	fileNumBuffers fileNumBuffers1
}

type fileReaders1 map[string]*bufio.Reader
// var fileReaders := make(fileReaders1)
type fileNumBuffers1 map[string][]int

func (t *T) openFiles(files []string) {
	t.fileReaders = make(map[string]*bufio.Reader)
	// fileReaders := make(fileReaders1)
	for _, path := range(files) {
		f, err := os.OpenFile(path, os.O_RDWR, 0755)
		if err != nil {
			log.Fatal(err)
		}
		// defer f.Close()
		reader := bufio.NewReader(f)
		// fileReaders[path] = make(*bufio.Reader)
		// 
		t.fileReaders[path] = reader
		t.fileNumBuffers[path] = make([]int, 0)
	}
}

func (t *T) fetchNumFromFiles() (int, error) {
	min_num_idx := ""
	var the_min_num int
	for path, data := range(t.fileNumBuffers) {
		if len(data) == 0 && t.fileReaders[path] != nil {

			// get from file
			reader := t.fileReaders[path]

			t.fileNumBuffers[path] = getNumsFromFile(reader, PER_GROUP_NUM, t.fileNumBuffers[path])

			if len(t.fileNumBuffers[path]) < PER_GROUP_NUM {
				// 取完了，就去掉
				delete(t.fileReaders, path)
			}

			// fileNumBuffers[path] = append(fileNumBuffers[path], get_nums)
			
			// append(fileNumBuffers[path], get_nums, fileNumBuffers[path])

			// log.Println(t.fileNumBuffers[path])
			// os.Exit(1)
		}

		data_len := len(t.fileNumBuffers[path])

		if data_len == 0 {
			// 文件排完序了，踢出
			delete(t.fileNumBuffers, path)
		} else {
			if min_num_idx == "" || t.fileNumBuffers[path][0] < the_min_num {
				min_num_idx = path
				the_min_num = t.fileNumBuffers[path][0]
			} 
			// fileNumBuffers[path][0]
			// fileNumBuffers[path] = fileNumBuffers[path][1:data_len]
		}
	}
	// log.Println(min_num_idx)

	if min_num_idx != "" {
		t.fileNumBuffers[min_num_idx] = t.fileNumBuffers[min_num_idx][1:]
		return the_min_num, nil
	}
	return the_min_num, errors.New("no new numbers")
}

func getNumsFromFile(reader *bufio.Reader, n int, nums []int) ([]int) {
	// var arr []int
	for {
		if n <= 0 {
			break
		}
		str, _, err := reader.ReadLine()
		if err != nil {
			// log.Println(err)
			break
		}

		tmp_int, _ := strconv.Atoi(string(str))
		// arr = append(arr, tmp_int)

		nums = append(nums, tmp_int)
		n--
	}
	// log.Println(nums)
	return nums
}


type ByCost []int

func (b ByCost) Len() int      { return len(b) }
func (b ByCost) Swap(i, j int) { b[i], b[j] = b[j], b[i] }

// 从低到高排序
func (b ByCost) Less(i, j int) bool { return b[i] < b[j] }

func sortFileLog(loginfos []int) {
	sort.Sort(ByCost(loginfos))
}