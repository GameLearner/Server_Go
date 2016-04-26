package Util
//IDGenerator generator uniq id
type idGenerator struct{
    id int;
}

var idGen idGenerator;

func init()  {
    idGen.id = 0;
}
//GetUniqID  return the uniq id in runtime
func GetUniqID() int {
    var id = idGen.id;
    idGen.id++;
    return id;
}