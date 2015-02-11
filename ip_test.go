package main

import (
        "testing"
        "net"
)

var test10ips  = [...]string{ "10.0.0.1", "10.16.0.1", "10.22.0.1", "10.30.0.1", "11.200.0.1", "10.255.255.255" }
var _,aclass,_ = net.ParseCIDR("10.0.0.0/8")
var test172ips = [...]string{ "172.16.0.1", "172.22.0.1", "172.30.0.1", "172.200.0.1", "172.255.255.255" }
var _,bclass,_ = net.ParseCIDR("172.16.0.0/12")
var test192ips = [...]string{ "192.168.0.1", "192.168.22.1", "192.169.30.1", "191.168.200.25", "192.168.255.255" }
var _,cclass,_ = net.ParseCIDR("192.168.0.0/16")


func Test10(t *testing.T) {
  for _,ip := range test10ips {
    var theIP = net.ParseIP(ip)
    if (aclass.Contains(theIP) != subnet10.MatchString(ip)) { 
      t.Errorf("Regex subnet10(%v) and net.Contains(%v) disagree about ip %v", subnet10.MatchString(ip), aclass.Contains(theIP), ip) 
    }
  }
}

func Test172(t *testing.T) {
  for _,ip := range test172ips {
    var theIP = net.ParseIP(ip)
    if (bclass.Contains(theIP) != subnet172.MatchString(ip)) { 
      t.Errorf("Regex subnet172(%v) and net.Contains(%v) disagree about ip %v", subnet172.MatchString(ip), bclass.Contains(theIP), ip) 
    }
  }
}

func Test192(t *testing.T) {
  for _,ip := range test192ips {
    var theIP = net.ParseIP(ip)
    if (cclass.Contains(theIP) != subnet192.MatchString(ip)) { 
      t.Errorf("Regex subnet192(%v) and net.Contains(%v) disagree about ip %v", subnet192.MatchString(ip), cclass.Contains(theIP), ip) 
    }
  }
}
