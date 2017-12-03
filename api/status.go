package main

import (
  "github.com/jinzhu/gorm"
)

type Status struct {
  gorm.Model

  Label string
  Color string
}

