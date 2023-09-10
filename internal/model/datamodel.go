package model

//// DataLoadAvg represents the load average of the system
////
////	Avg1 - The average processor workload of the last minute
////	Avg2 - The average processor workload of the last 5 minutes
////	Avg3 - The average processor workload of the last 15 minutes
//type DataLoadAvg struct {
//	Avg1 float64
//	Avg2 float64
//	Avg3 float64
//}
//
//// DataCpuAvgStats represents the average cpu load
////
////	User, System, Idle
//type DataCpuAvgStats struct {
//	User   float64
//	System float64
//	Idle   float64
//}
//
//// DataDisksLoad represents the  Disks Load
////
////	Device  -like nvme, sda,  etc
////	Tps - transfers per second
////	Kbps - KB/s (kilobytes (read+write) per second
//type DataDisksLoad struct {
//	Device string
//	Tps    int64
//	Kbps   float64
//}
//
//// DataDisksUsage represents the  Disks usage information
////
////	MountPoint   - mount point
////	FileSystem   - file system(ext4, etc)
////	Inode        - usage inodes
////	InodePercent - usage inodes percent
////	Mb           - usage Mbytes
////	MbPercent    - usage Mbytes percent
//type DataDisksUsage struct {
//	MountPoint   string
//	FileSystem   string
//	Inode        int64
//	InodePercent float64
//	Mb           float64
//	MbPercent    float64
//}
//
//// DataTopNetworkProto represents the  top talkers network by proto
////
////	Proto -        Protocol (TCP, UDP..)
////	Bytes  -       bytes
////	BytesPercent - bytes percent
//type DataTopNetworkProto struct {
//	Proto        string
//	Bytes        int64
//	BytesPercent float64
//}
//
//// DataTopNetworkTraffic represents the  top talkers network by Source/Dest
////
////	Source      - source:port
////	Destination - destination:port
////	Proto       - protocol
////	Bps			- bytes per second
//type DataTopNetworkTraffic struct {
//	Source      string
//	Destination string
//	Proto       string
//	Bps         int64
//}
//
//// DataNetworkListen represents the stat of listen connections
////
////	Proto   - protocol
////	Command - command
////	Pid     - pid
////	User    - user
////	Port    - port
//type DataNetworkListen struct {
//	Proto   string
//	Command string
//	Pid     int64
//	User    string
//	Port    int64
//}
//
//// DataNetworkStates represents the stat of listen connections
////
////	State - ESTAB, FIN_WAIT, SYN_RCV, etc
////	Number - number of states connections
//type DataNetworkStates struct {
//	State  string
//	Number int64
//}
