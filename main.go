package main

func main() {
		dal, _ := newDal("test.db")

		p := dal.allocateEmptyPage()
		p.num = dal.getNextPage()

		copy(p.data[:], []byte("Hello, World!"))

		_ = dal.writePage(p)
		_, _ = dal.writeFreelist()

		_ = dal.close()

		dal, _ = newDal("test.db")

		p = dal.allocateEmptyPage()
		p.num = dal.getNextPage()

		copy(p.data[:], []byte("Hello, World 2 !"))

		_ = dal.writePage(p)

		pageNum := dal.getNextPage()
		dal.releasePage(pageNum)

		_, _ = dal.writeFreelist()
}
