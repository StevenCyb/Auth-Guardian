package mocked

import (
	"auth-guardian/config"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	ber "gopkg.in/asn1-ber.v1"
)

// RunMockLDAPIDP runs a mocked LDAP IDP
func RunMockLDAPIDP() {
	listener, err := net.Listen("tcp", ":3002")
	if err != nil {
		log.Panic(err)
	}

	defer listener.Close()
	for {
		// Wait for connection
		connection, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}

		// Keep the connection for a minute
		connection.SetDeadline(time.Now().Add(time.Minute))

		go func(c net.Conn) {
			for {
				// Read and parse asn data
				p, err := ber.ReadPacket(c)
				if err != nil {
					break
				}

				// Get message id
				mID, err := strconv.ParseInt(fmt.Sprint(p.Children[0].Value), 10, 64)
				if err != nil {
					continue
				}

				if mID == 1 {
					// Initial bind request

					// Check bind password
					if string(p.Children[1].Children[2].Data.Bytes()) != config.DirectoryServerBindPassword {
						log.Fatal("Wrong bind password")
					}

					// Check bind dn
					if string(p.Children[1].Children[1].ByteValue) != config.DirectoryServerBindDN {
						log.Fatal("Wrong bind dn")
					}

					// Send positive bind response and send back
					shortResponse(c, mID, 1, 0)
				} else if mID == 2 {
					// Search request

					// Check base dn
					if string(p.Children[1].Children[0].ByteValue) != config.DirectoryServerBaseDN {
						log.Fatal("Wrong base dn")
					}

					// Check username
					if string(p.Children[1].Children[6].Children[1].ByteValue) != "san" {
						// Build search result done response and send back since username is not correct
						shortResponse(c, mID, 5, 0)
					}

					replypacket1 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "LDAP Response")
					replypacket1.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, mID, "MessageId"))

					// Build search result entry response
					uid := string(p.Children[1].Children[6].Children[1].ByteValue)

					// Response id = 4 = search result entry
					searchResult := ber.Encode(ber.ClassApplication, ber.TypeConstructed, ber.Tag(4), nil, "Response")
					searchResult.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "uid="+uid+",dc=example,dc=com", "Object Name"))

					// Define response object
					object := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Object")

					// Sub object for object class
					subObject1 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Sub Object 1")
					subObject1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "objectClass", "SO1 Attribute Name"))
					subObject1_1 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "Sub Object 1_1")
					subObject1_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "princess", "SO1_1 Attribute"))
					subObject1_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "mononoke", "SO1_1 Attribute"))
					subObject1_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "wolf", "SO1_1 Attribute"))
					subObject1_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "women", "SO1_1 Attribute"))
					subObject1.AppendChild(subObject1_1)
					object.AppendChild(subObject1)

					// Sub object for cn
					subObject2 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Sub Object 2")
					subObject2.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "cn", "SO2 Attribute Name"))
					subObject2_1 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "Sub Object 2_1")
					subObject2_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "Princess San", "SO2_1 Attribute"))
					subObject2.AppendChild(subObject2_1)
					object.AppendChild(subObject2)

					// Sub object for sn
					subObject3 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Sub Object 3")
					subObject3.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "sn", "SO3 Attribute Name"))
					subObject3_1 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "Sub Object 3_1")
					subObject3_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "San", "SO3_1 Attribute"))
					subObject3.AppendChild(subObject3_1)
					object.AppendChild(subObject3)

					// Sub object for uid
					subObject4 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Sub Object 4")
					subObject4.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "uid", "SO4 Attribute Name"))
					subObject4_1 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "Sub Object 4_1")
					subObject4_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "san", "SO4_1 Attribute"))
					subObject4.AppendChild(subObject4_1)
					object.AppendChild(subObject4)

					// Sub object for mail
					subObject5 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Sub Object 5")
					subObject5.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "mail", "SO5 Attribute Name"))
					subObject5_1 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "Sub Object 5_1")
					subObject5_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "san@local-mock.com", "SO5_1 Attribute"))
					subObject5.AppendChild(subObject5_1)
					object.AppendChild(subObject5)

					// Sub object for uid no
					subObject6 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Sub Object 6")
					subObject6.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "uidNumber", "SO6 Attribute Name"))
					subObject6_1 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "Sub Object 6_1")
					subObject6_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "1997", "SO6_1 Attribute"))
					subObject6.AppendChild(subObject6_1)
					object.AppendChild(subObject6)

					// Sub object for gid no
					subObject7 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Sub Object 7")
					subObject7.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "gidNumber", "SO7 Attribute Name"))
					subObject7_1 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "Sub Object 7_1")
					subObject7_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "2001", "SO6_7 Attribute"))
					subObject7.AppendChild(subObject7_1)
					object.AppendChild(subObject7)

					// Sub object for home dir
					subObject8 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "Sub Object 8")
					subObject8.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "homeDirectory", "SO8 Attribute Name"))
					subObject8_1 := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "Sub Object 8_1")
					subObject8_1.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "forest", "SO6_8 Attribute"))
					subObject8.AppendChild(subObject8_1)
					object.AppendChild(subObject8)

					searchResult.AppendChild(object)
					replypacket1.AppendChild(searchResult)

					// Send search result entry back
					send(c, replypacket1)

					// Build search result done response and send back
					shortResponse(c, mID, 5, 0)
				} else if mID == 3 {
					// Validate credentials

					// Get username and password
					username := string(p.Children[1].Children[1].ByteValue)
					password := string(p.Children[1].Children[2].Data.Bytes())

					if username == "uid=san,dc=example,dc=com" && password == "Wolf" {
						// Build positive bind response and send back
						shortResponse(c, mID, 1, 0)
					} else {
						// Build negative bind response and send back
						shortResponse(c, mID, 1, 40)
					}
					break
				}
			}

			// Shut down the connection
			c.Close()
		}(connection)
	}
}

func send(c net.Conn, replypacket *ber.Packet) {
	for _, r := range []*ber.Packet{replypacket} {
		fmt.Println("Package:")
		c.Write(r.Bytes())
	}
}

func shortResponse(c net.Conn, mID int64, tag int64, resultCode int64) {
	// Build response
	replypacket := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "LDAP Response")
	replypacket.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, mID, "MessageId"))
	// Response id https://ldapwiki.com/wiki/LDAP%20Message
	bindResult := ber.Encode(ber.ClassApplication, ber.TypeConstructed, ber.Tag(tag), nil, "Response")
	bindResult.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagEnumerated, resultCode, "Result Code"))
	bindResult.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "", "Unused"))
	bindResult.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, "", "Unused"))
	replypacket.AppendChild(bindResult)

	// Send response
	send(c, replypacket)
}
