## Potok Orm
[![Build Status](https://travis-ci.org/adtechpotok/go-orm.svg?branch=master)](https://travis-ci.org/adtechpotok/go-orm)
-
2 type are supported: sliced and base with stringed or int id.
Base will return interface or nil.
Sliced will return slice.
## Example base struct
```
 type UserInfo struct {
 	Id           int    `gorm:"column:oid"`
 	Status       string `gorm:"column:status"`
 }
 
 func (m UserInfo) IsActive() bool {
 	return m.Status == orm.StatusActive
 }
 
 func (m UserInfo) GetId() int {
 	return m.Id
 }
 
 type UserInfoRep struct {
 	orm.BaseDbModel
 }
 
 func (m *UserInfoRep) InitialSelect(db *gorm.DB) {
 	var r []UserInfo
 	db.Where("status = ?", orm.StatusActive).Find(&r)
 
 	for _, item := range r {
 		m.AddToCache(item)
 	}
 
 }
```
or geting data in code
 ```$xslt
	var userInfoRep = rep.UserInfoRep{}
	userInfoRep.InitialSelect(db)
	userInfoVal := userInfoCache.FindInCache(1)
	if userInfoVal != nil {
		userInfo := userInfoVal.(UserInfo)
	}
```