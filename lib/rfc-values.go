package lib

/*
3.2.2. TYPE values

TYPE fields are used in resource records.  Note that these types are a
subset of QTYPEs.

TYPE            value and meaning

A               1 a host address

NS              2 an authoritative name server

MD              3 a mail destination (Obsolete - use MX)

MF              4 a mail forwarder (Obsolete - use MX)

CNAME           5 the canonical name for an alias

SOA             6 marks the start of a zone of authority

MB              7 a mailbox domain name (EXPERIMENTAL)

MG              8 a mail group member (EXPERIMENTAL)

MR              9 a mail rename domain name (EXPERIMENTAL)

NULL            10 a null RR (EXPERIMENTAL)

WKS             11 a well known service description

PTR             12 a domain name pointer

HINFO           13 host information

MINFO           14 mailbox or mail list information

MX              15 mail exchange

TXT             16 text strings

*/

type TYPE uint8

const (
	TypeA     TYPE = iota + 1 // 1 a host address
	TypeNS                    // 2 an authoritative name server
	_                         // 3 a mail destination (Obsolete - use MX)
	_                         // 4 a mail forwarder (Obsolete - use MX)
	TypeCNAME                 // 5 the canonical name for an alias
	TypeSOA                   // 6 marks the start of a zone of authority
	TypeMB                    // 7 a mailbox domain name (EXPERIMENTAL)
	TypeMG                    // 8 a mail group member (EXPERIMENTAL)
	TypeMR                    // 9 a mail rename domain name (EXPERIMENTAL)
	TypeNULL                  // 10 a null RR (EXPERIMENTAL)
	TypeWKS                   // 11 a well known service description
	TypePTR                   // 12 a domain name pointer
	TypeHINFO                 // 13 host information
	TypeMINFO                 // 14 mailbox or mail list information
	TypeMX                    // 15 mail exchange
	TypeTXT                   // 16 text strings
)

/*
3.2.3. QTYPE values

QTYPE fields appear in the question part of a query.  QTYPES are a
superset of TYPEs, hence all TYPEs are valid QTYPEs.  In addition, the
following QTYPEs are defined:

AXFR            252 A request for a transfer of an entire zone

MAILB           253 A request for mailbox-related records (MB, MG or MR)

MAILA           254 A request for mail agent RRs (Obsolete - see MX)

*               255 A request for all records

*/

type QTYPE uint8

const (
	TypeAXFR = iota + 252
	TypeMAILB
	TypeMAILA
	TypeALL
)

/*
3.2.4. CLASS values

CLASS fields appear in resource records.  The following CLASS mnemonics
and values are defined:

IN              1 the Internet

CS              2 the CSNET class (Obsolete - used only for examples in
                some obsolete RFCs)

CH              3 the CHAOS class

HS              4 Hesiod [Dyer 87]

*/

type CLASS uint8

type QCLASS CLASS

const (
	ClassIN = iota + 1 // 1 the Internet
	ClassCS            // 2 the CSNET class (Obsolete - used only for examples in
	// some obsolete RFCs)
	ClassCH // 3 the CHAOS class
	ClassHS // 4 Hesiod [Dyer 87]
)

/*

3.2.5. QCLASS values

QCLASS fields appear in the question section of a query.  QCLASS values
are a superset of CLASS values; every CLASS is a valid QCLASS.  In
addition to CLASS values, the following QCLASSes are defined:

*               255 any class
*/

const (
	ClassANY = 255
)

/*
3.3.1. CNAME RDATA format

    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    /                     CNAME                     /
    /                                               /
    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

where:

CNAME           A <domain-name> which specifies the canonical or primary
                name for the owner.  The owner name is an alias.

CNAME RRs cause no additional section processing, but name servers may
choose to restart the query at the canonical name in certain cases.  See
the description of name server logic in [RFC-1034] for details.
*/

type CNAME string
