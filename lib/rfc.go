package lib

/*
https://tools.ietf.org/html/rfc1035

3.2.1. Format

All RRs have the same top level format shown below:

                                    1  1  1  1  1  1
      0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                                               |
    /                                               /
    /                      NAME                     /
    |                                               |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                      TYPE                     |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                     CLASS                     |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                      TTL                      |
    |                                               |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    |                   RDLENGTH                    |
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
    /                     RDATA                     /
    /                                               /
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+


where:

NAME            an owner name, i.e., the name of the node to which this
                resource record pertains.

TYPE            two octets containing one of the RR TYPE codes.

CLASS           two octets containing one of the RR CLASS codes.

TTL             a 32 bit signed integer that specifies the time interval
                that the resource record may be cached before the source
                of the information should again be consulted.  Zero
                values are interpreted to mean that the RR can only be
                used for the transaction in progress, and should not be
                cached.  For example, SOA records are always distributed
                with a zero TTL to prohibit caching.  Zero values can
                also be used for extremely volatile data.

RDLENGTH        an unsigned 16 bit integer that specifies the length in
				octets of the RDATA field.

RDATA           a variable length string of octets that describes the
                resource.  The format of this information varies
				according to the TYPE and CLASS of the resource record.

*/

type RR struct {
	Name     string
	Type     TYPE
	Class    CLASS
	TTL      int32
	RDLENGTH uint16
	RDATA    string
}

/*
4. MESSAGES

4.1. Format

All communications inside of the domain protocol are carried in a single
format called a message.  The top level format of message is divided
into 5 sections (some of which are empty in certain cases) shown below:

    +---------------------+
    |        Header       |
    +---------------------+
    |       Question      | the question for the name server
    +---------------------+
    |        Answer       | RRs answering the question
    +---------------------+
    |      Authority      | RRs pointing toward an authority
    +---------------------+
    |      Additional     | RRs holding additional information
    +---------------------+

The header section is always present.  The header includes fields that
specify which of the remaining sections are present, and also specify
whether the message is a query or a response, a standard query or some
other opcode, etc.

The names of the sections after the header are derived from their use in
standard queries.  The question section contains fields that describe a
question to a name server.  These fields are a query type (QTYPE), a
query class (QCLASS), and a query domain name (QNAME).  The last three
sections have the same format: a possibly empty list of concatenated
resource records (RRs).  The answer section contains RRs that answer the
question; the authority section contains RRs that point toward an
authoritative name server; the additional records section contains RRs
which relate to the query, but are not strictly answers for the
question.
*/

type MESSAGES struct {
}

type HEADER struct {
	ID      [16]byte
	QR      byte
	OPCODE  [4]byte
	AA      byte
	TC      byte
	RD      byte
	RA      byte
	Z       [3]byte
	RCODE   [4]byte
	QDCOUNT [16]byte
	ANCOUNT [16]byte
	NSCOUNT [16]byte
	ARCOUNT [16]byte
}

type QUESTION struct {
}
