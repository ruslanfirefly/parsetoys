package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)
const (
	DBMNG string = "parsetoys"
	COLLINK string = "visitedlinks"
	COLTOYS string = "toysinfo"
)

func createMongoSession() *mgo.Session{
	session, err := mgo.Dial(mongoURL)
	error_log(err)
	return session
}

func dropOldLinks(){
	s := createMongoSession()
	defer s.Close()
	listCol, err := s.DB(DBMNG).C(COLLINK).Count()
	error_log(err)
	if(listCol != 0){
		err = s.DB(DBMNG).C(COLLINK).DropCollection()
		error_log(err)
	}

}

func inBaseDonor(donor Donor) bool{
	s := createMongoSession()
	defer s.Close()
	listRes, err := s.DB(DBMNG).C(COLLINK).Find(donor).Count()
	error_log(err)
	if (listRes == 0){
		return false
	}else{
		return true
	}
}
func inBaseTovar(tovar Tovar) bool{
	s := createMongoSession()
	defer s.Close()
	listRes, err := s.DB(DBMNG).C(COLTOYS).Find(bson.M{"site":tovar.Site,"art":tovar.Art}).Count()
	error_log(err)
	if (listRes == 0){
		return false
	}else{
		return true
	}
}
func addVisitedLinks(donor Donor){
	s := createMongoSession()
	defer s.Close()
	err := s.DB(DBMNG).C(COLLINK).Insert(donor)
	error_log(err)
}

func saveTovar(tovar Tovar){
	var err error
	s := createMongoSession()
	defer s.Close()
	if(inBaseTovar(tovar)){
		err = s.DB(DBMNG).C(COLTOYS).Update(bson.M{"site":tovar.Site,"art":tovar.Art},tovar)
	}else{
		err = s.DB(DBMNG).C(COLTOYS).Insert(tovar)
	}
	error_log(err)
	wg.Done()
}
