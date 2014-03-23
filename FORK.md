This project is a fork from github.com/dmotylev/goproperties and adds or changes the following:

- (add) setting for linebuffer size (properties.ReaderLineBufferSize = 1024)
- (change) use that setting
- (add) get a subset of properties based on regular expression on key
   
   mongodProperties := properties.SelectProperties("mongod.*")

https://github.com/emicklei