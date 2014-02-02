package main

import (
	"labix.org/v2/mgo"
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
func isNotInBase(donor Donor) bool{
	s := createMongoSession()
	defer s.Close()
	listRes, err := s.DB(DBMNG).C(COLLINK).Find(donor).Count()
	error_log(err)
	if (listRes == 0){
		return true
	}else{
		return false
	}
}
func addVisitedLinks(donor Donor){
	s := createMongoSession()
	defer s.Close()
	err := s.DB(DBMNG).C(COLLINK).Insert(donor)
	error_log(err)
}
