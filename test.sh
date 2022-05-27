./publisher -ch=foo -data=testData/model1.json
./publisher -ch=foo -data=testData/model2.json
./publisher -ch=foo -data=testData/model3.json
# try publish files with errors
./publisher -ch=foo -data=testData/modelError.json
./publisher -ch=foo -data=testData/err.sh
./publisher -ch=foo -data=testData/emptyError.json
# try publish model that already exists
./publisher -ch=foo -data=testData/model2.json