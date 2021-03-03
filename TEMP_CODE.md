err = ioutil.WriteFile("./" + strconv.Itoa(task.Id) + "-" + strconv.Itoa(i) + ".jpg", *buf, 0644)
if err != nil {
    fmt.Println("ERR.JobHandler.Run.Screenshot.2", err)
}
