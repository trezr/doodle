package main

import (
  "github.com/jinzhu/gorm"
)

type Participation struct {
  gorm.Model

  Name string
  Value bool

  EventID int
}
