/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package redis

import (
	"time"

	"github.com/blend/go-sdk/ex"
)

// Defaults
const (
	DefaultNetwork		= "tcp"
	DefaultAddr		= "127.0.0.1:6379"
	DefaultTimeout		= 5 * time.Second
	DefaultConnectTimeout	= time.Second
)

// Errors
const (
	ErrPingFailed ex.Class = "radix ping failed"
)

// Key Operations
const (
	// OpCOPY copies the value stored at the source key to the destination key.
	//
	// By default, the destination key is created in the logical database used by the connection. The DB option allows specifying an alternative logical database index for the destination key.
	//
	// The command returns an error when the destination key already exists. The REPLACE option removes the destination key before copying the value to it.
	//
	// Usage: COPY source destination [DB destination-db] [REPLACE]
	//
	// Return value is an integer reply, specifically:
	//
	//    1 if source was copied.
	//    0 if source was not copied.
	//
	OpCOPY	= "COPY"

	// OpDel removes the specified keys. A key is ignored if it does not exist.
	//
	// Usage: DEL key [key ...]
	//
	// Return value is an integer reply: The number of keys that were removed.
	OpDEL	= "DEL"

	// OpDUMP serialies the value stored at key in a Redis-specific format and return it to the user. The returned value can be synthesized back into a Redis key using the RESTORE command.
	//
	// The serialized value does NOT contain expire information. In order to capture the time to live of the current value the PTTL command should be used.
	//
	// Usage: DUMP key
	//
	// Return value is a bulk string reply: the serialized value.
	OpDUMP	= "DUMP"

	// OpEXISTS returns if key exists.
	//
	// The user should be aware that if the same existing key is mentioned in the arguments multiple times, it will be counted multiple times. So if somekey exists, EXISTS somekey somekey will return 2.
	//
	// Usage: EXISTS key [key ...]
	//
	// Returns an integer reply, specifically:
	//
	//    1 if the key exists.
	//    0 if the key does not exist.
	//
	OpEXISTS	= "EXISTS"

	// OpExpire sets a timeout on key. After the timeout has expired, the key will automatically be deleted. A key with an associated timeout is often said to be volatile in Redis terminology.
	//
	// The timeout will only be cleared by commands that delete or overwrite the contents of the key, including DEL, SET, GETSET and all the *STORE
	//
	// Usage: EXPIRE key seconds
	//
	// Return value is an integer reply, specifically:
	//
	//    1 if the timeout was set.
	//    0 if key does not exist.
	//
	OpEXPIRE	= "EXPIRE"
	OpEXPIREAT	= "EXPIREAT"
	OpKEYS		= "KEYS"
	OpMIGRATE	= "MIGRATE"
	OpMOVE		= "MOVE"
	OpOBJECT	= "OBJECT"
	OpPERSIST	= "PERSIST"
	OpPEXPIRE	= "PEXPIRE"
	OpPEXPIREAT	= "PEXPIREAT"
	OpPTTL		= "PTTL"
	OpRANDOMKEY	= "RANDOMKEY"
	OpRENAME	= "RENAME"
	OpRENAMENX	= "RENAMENX"
	OpRESTORE	= "RESTORE"
	OpSORT		= "SORT"
	OpTOUCH		= "TOUCH"
	OpTTL		= "TTL"
	OpTYPE		= "TYPE"
	OpUNLINK	= "UNLINK"
	OpWAIT		= "WAIT"
	OpSCAN		= "SCAN"
)

// Hash Operations
const (
	OpHDEL		= "HDEL"
	OpHEXISTS	= "HEXISTS"
	OpHGET		= "HGET"
	OpHGETALL	= "HGETALL"
	OpHINCRBY	= "HINCRBY"
	OpHINCRBYFLOAT	= "HINCRBYFLOAT"
	OpHKEYS		= "HKEYS"
	OpHLEN		= "HLEN"
	OpHMGET		= "HMGET"
	OpHMSET		= "HMSET"
	OpHSET		= "HSET"
	OpHSETNX	= "HSETNX"
	OpHRANDFIELD	= "HRANDFIELD"
	OpHSTRLEN	= "HSTRLEN"
	OpHVALS		= "HVALS"
	OpHSCAN		= "HSCAN"
)

// Set Operations
const (
	// OpSADD adds the specified members to the set stored at key. Specified members that are already a member of this set are ignored. If key does not exist, a new set is created before adding the specified members.
	//
	// Usage: SADD key member [member ...]
	//
	// An error is returned when the value stored at key is not a set.
	OpSADD	= "SADD"

	// OpSCARD returns the set cardinality (number of elements) of the set stored at key.
	//
	// Usage: SCARD key
	//
	// Return value is an integer reply: the cardinality (number of elements) of the set, or 0 if key does not exist.
	OpSCARD	= "SCARD"

	// OpSDIFF returns the members of the set resulting from the difference between the first set and all the successive sets.
	//
	// Keys that do not exist are considered to be empty sets.
	//
	// Usage: SDIFF key [key ...]
	//
	// Return value is an array reply: list with members of the resulting set.
	OpSDIFF	= "SDIFF"

	// OpSDIFFSTORE is command is equal to SDIFF, but instead of returning the resulting set, it is stored in destination.
	//
	// If destination already exists, it is overwritten.
	//
	// Usage: SDIFFSTORE destination key [key ...]
	//
	// Return value is an integer reply: the number of elements in the resulting set.
	OpSDIFFSTORE	= "SDIFFSTORE"

	// OpSINTER returns the members of the set resulting from the intersection of all the given sets.
	//
	// Keys that do not exist are considered to be empty sets. With one of the keys being an empty set, the resulting set is also empty (since set intersection with an empty set always results in an empty set).
	//
	// Usage: SINTER key [key ...]
	//
	// Return value is an array reply: list with members of the resulting set.
	OpSINTER	= "SINTER"

	// OpSINTERSTORE is equal to SINTER, but instead of returning the resulting set, it is stored in destination. If destination already exists, it is overwritten.
	//
	// Usage: SINTERSTORE destination key [key ...]
	//
	// Return value is an integer reply: the number of elements in the resulting set.
	OpSINTERSTORE	= "SINTERSTORE"

	// OpSISMEMBER returns if member is a member of the set stored at key.
	//
	// Usage: SISMEMBER key member
	//
	// Return value is an integer reply, specifically:
	//
	//    1 if the element is a member of the set.
	//    0 if the element is not a member of the set, or if key does not exist.
	//
	OpSISMEMBER	= "SISMEMBER"

	// OpSMISMEMBER returns whether each member is a member of the set stored at key.
	//
	// For every member, 1 is returned if the value is a member of the set, or 0 if the element is not a member of the set or if key does not exist.
	//
	// Usage: SMISMEMBER key member [member ...]
	//
	// Return value is an array reply: list representing the membership of the given elements, in the same order as they are requested.
	OpSMISMEMBER	= "SMISMEMBER"

	// OpSMEMBERS returns all the members of the set value stored at key. This has the same effect as running SINTER with one argument key.
	//
	// Usage: SMEMBERS key
	//
	// Return value is an array reply: all elements of the set.
	OpSMEMBERS	= "SMEMBERS"

	// OpSMOVE moves a member from the set at source to the set at destination. This operation is atomic. In every given moment the element will appear to be a member of source or destination for other clients.
	//
	// If the source set does not exist or does not contain the specified element, no operation is performed and 0 is returned. Otherwise, the element is removed from the source set and added to the destination set. When the specified element already exists in the destination set, it is only removed from the source set.
	//
	// An error is returned if source or destination does not hold a set value.
	//
	// Usage: SMOVE source destination member
	//
	// Return value is an integer reply, specifically:
	//
	//    1 if the element is moved.
	//    0 if the element is not a member of source and no operation was performed.
	//
	OpSMOVE	= "SMOVE"

	// OpSPOP removes and returns one or more random members from the set value store at key.
	//
	// This operation is similar to SRANDMEMBER, that returns one or more random elements from a set but does not remove it.
	//
	// By default, the command pops a single member from the set. When provided with the optional count argument, the reply will consist of up to count members, depending on the set's cardinality.
	//
	// Usage: SPOP key [count]
	//
	// Return value when called without the count argument:
	//
	// Bulk string reply: the removed member, or nil when key does not exist.
	//
	// When called with the count argument:
	//
	// Array reply: the removed members, or nil when key does not exist.
	OpSPOP	= "SPOP"

	// OpSRANDMEMBER returns a random element from the set value stored at key, or given a count, returns an array of distinct elements. The array's length is either count or the set's cardinality (SCARD), whichever is lower.
	//
	// Usage: SRANDMEMBER key [count]
	//
	// Return value is bulk a string reply: without the additional count argument, the command returns a Bulk Reply with the randomly selected element, or nil when key does not exist.
	//
	// Array reply: when the additional count argument is passed, the command returns an array of elements, or an empty array when key does not exist.
	OpSRANDMEMBER	= "SRANDMEMBER"

	// OpSREM removes the specified members from the set stored at key. Specified members that are not a member of this set are ignored. If key does not exist, it is treated as an empty set and this command returns 0.
	//
	// An error is returned when the value stored at key is not a set.
	//
	// Usage: SREM key member [member ...]
	//
	// Return value is an integer reply: the number of members that were removed from the set, not including non existing members.
	OpSREM	= "SREM"

	// OpSSCAN is used to incrementally iterate over a collection of elements, specifically the set type.
	//
	// Basic usage: SCAN is a cursor based iterator. This means that at every call of the command, the server returns an updated cursor that the user needs to use as the cursor argument in the next call.
	//
	// An iteration starts when the cursor is set to 0, and terminates when the cursor returned by the server is 0.
	OpSSCAN	= "SSCAN"

	// OpSUNION returns the members of the set resulting from the union of all the given sets.
	//
	// Usage: SUNION key [key ...]
	//
	// Return value is an array reply: list with members of the resulting set.
	OpSUNION	= "SUNION"

	// OpSUNIONSTORE is equal to SUNION, but instead of returning the resulting set, it is stored in destination. If destination already exists, it is overwritten.
	//
	// Usage: SUNIONSTORE destination key [key ...]
	//
	// Return value is an integer reply: the number of elements in the resulting set.
	OpSUNIONSTORE	= "SUNIONSTORE"
)

// String Operations
const (
	OpAPPEND	= "APPEND"
	OpBITCOUNT	= "BITCOUNT"
	OpBITFIELD	= "BITFIELD"
	OpBITOP		= "BITOP"
	OpBITPOS	= "BITPOS"
	OpDECR		= "DECR"
	OpDECRBY	= "DECRBY"
	// OpGet gets the value of key. If the key does not exist the special value nil is returned. An error is returned if the value stored at key is not a string, because GET only handles string values.
	//
	// Usage: GET key
	//
	// Return value is a bulk string reply: the value of key, or nil when key does not exist.
	OpGET		= "GET"
	OpGETBIT	= "GETBIT"
	OpGETDEL	= "GETDEL"
	// OpGETEX gets the value of key and optionally set its expiration. GETEX is similar to GET, but is a write command with additional options.
	//
	// Usagee: GETEX key [EX seconds|PX milliseconds|EXAT timestamp|PXAT milliseconds-timestamp|PERSIST]
	//
	// Return value is a bulk string reply: the value of key, or nil when key does not exist.
	OpGETEX		= "GETEX"
	OpGETRANGE	= "GETRANGE"
	OpGETSET	= "GETSET"
	OpINCR		= "INCR"
	OpINCRBY	= "INCRBY"
	OpINCRBYFLOAT	= "INCRBYFLOAT"
	OpMGET		= "MGET"
	OpMSET		= "MSET"
	OpMSETNX	= "MSETNX"
	OpPSETEX	= "PSETEX"
	OpSET		= "SET"
	OpSETBIT	= "SETBIT"

	// OpSETEX sets a key to hold the string value and set key to timeout after a given number of seconds.
	//
	// Usage: SETEX key seconds value
	//
	// Return value is a simple string reply.
	OpSETEX		= "SETEX"
	OpSETNX		= "SETNX"
	OpSETRANGE	= "SETRANGE"
	OpSTRALGO	= "STRALGO"
	OpSTRLEN	= "STRLEN"
)

// Connection Operations
const (
	OpAUTH			= "AUTH"
	OpCLIENTCACHING		= "CLIENT CACHING"
	OpCLIENTID		= "CLIENT ID"
	OpCLIENTINFO		= "CLIENT INFO"
	OpCLIENTKILL		= "CLIENT KILL"
	OpCLIENTLIST		= "CLIENT LIST"
	OpCLIENTGETNAME		= "CLIENT GETNAME"
	OpCLIENTGETREDIR	= "CLIENT GETREDIR"
	OpCLIENTUNPAUSE		= "CLIENT UNPAUSE"
	OpCLIENTPAUSE		= "CLIENT PAUSE"
	OpCLIENTREPLY		= "CLIENT REPLY"
	OpCLIENTSETNAME		= "CLIENT SETNAME"
	OpCLIENTTRACKING	= "CLIENT TRACKING"
	OpCLIENTTRACKINGINFO	= "CLIENT TRACKINGINFO"
	OpCLIENTUNBLOCK		= "CLIENT UNBLOCK"
	OpECHO			= "ECHO"
	OpHELLO			= "HELLO"
	OpPING			= "PING"
	OpQUIT			= "QUIT"
	OpRESET			= "RESET"
	OpSELECT		= "SELECT"
)
