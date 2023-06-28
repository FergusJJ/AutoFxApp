package fix

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func (session *FxSession) sendMessage(message string, user FxUser) *FxResponse {
	requestBody := []byte(message)
	resp, err := readAndWrite(session.Connection, requestBody)
	if err != nil {
		log.Println(err)
		return &FxResponse{body: []byte{}, err: err}
	}
	return &FxResponse{body: resp, err: nil}

}

// create connection, then need to use connection to send message, read message back then close connection
// maybe look at caching the connection in order to speed up?
func CreateConnection(hostName string, port int) (conn *tls.Conn, err error) {
	tcpAddress, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", hostName, port))
	if err != nil {
		return conn, err
	}
	tlsConfig := getTlsConfig(hostName)
	conn, err = tls.Dial(
		"tcp",
		tcpAddress.String(),
		tlsConfig,
	)
	if err != nil {
		return conn, err
	}
	if err != conn.Handshake() {
		return conn, err
	}
	return conn, err
}

func getTlsConfig(hostName string) *tls.Config {
	config := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         hostName,
		//RootCAs:            x509.NewCertPool(),
	}
	return config
}

// looks like will need allow to read more bytes
// otherwise, it will return the previous message in the
func readAndWrite(conn *tls.Conn, reqBytes []byte) (buf []byte, err error) {

	_, err = conn.Write(reqBytes)
	if err != nil {
		return buf, err
	}
	var tmp = make([]byte, 1024)
	// var tmp = make([]byte, 200000)
	var totalBytesRead = 0
	for {
		l, err := conn.Read(tmp)
		if err != nil {
			return buf, err
		}
		totalBytesRead = totalBytesRead + l
		buf = append(buf, tmp...)
		if l == 1024 {
			continue
		}
		break
	}
	buf = append(buf, []byte{1, 2, 2, 2, 2, 2, 2}...)
	// buf = append(buf, []byte{}...)
	return buf, err
}

/*

[8=FIX.4.4 9=83406 35=y 34=2 49=CServer 50=TRADE 52=20230507-00:53:01.270 56=demo.ctrader.3697899 57=TRADE 320=Sxo2Xlb1jzJC 322=responce:Sxo2Xlb1jzJC 560=0 146=2462 55=20480 1007=CAMECO CORP 1008=2 55=1 1007=EURUSD 1008=5 55=20481 1007=GRUMA SAB-ADR (NO SHORT-SELLING) 1008=2 55=2 1007=GBPUSD 1008=5 55=20482 1007=MULTIMEDIA GAMES 1008=2 55=3 1007=EURJPY 1008=3 55=20483 1007=ANHEUSER BUSCH 1008=2 55=4 1007=USDJPY 1008=3 55=20484 1007=MELLANOX TECH 1008=2 55=5 1007=AUDUSD 1008=5 55=20485 1007=GOOGLE CLASS A 1008=2 55=6 1007=USDCHF 1008=5 55=20486 1007=BROADCOM CORP 1008=2 55=7 1007=GBPJPY 1008=3 55=20487 1007=FIRST ENERGY 1008=2 55=8 1007=USDCAD 1008=5 55=20488 1007=BANK OF AMERICA 1008=2 55=9 1007=EURGBP 1008=5 55=20489 1007=ADTRAN INC 1008=2 55=10 1007=EURCHF 1008=5 55=20490 1007=NATIONAL RETAIL 1008=2 55=11 1007=AUDJPY 1008=3 55=20491 1007=NIKE INC CL B 1008=2 55=12 1007=NZDUSD 1008=5 55=20492 1007=GREEN MOUNTAIN COFFEE 1008=2 55=13 1007=CHFJPY 1008=3 55=20493 1007=EQ RESIDENT 1008=2 55=14 1007=EURAUD 1008=4 55=20494 1007=BRITISH AMER TOB 1008=2 55=15 1007=CADJPY 1008=3 55=20495 1007=TRACTOR SUPPLY 1008=2 55=16 1007=GBPAUD 1008=5 55=20496 1007=MARKETAXESS HLDGS 1008=2 55=17 1007=EURCAD 1008=5 55=20497 1007=ACTIVISION INC NEW 1008=2 55=18 1007=AUDCAD 1008=5 55=20498 1007=CUBIST PHARMACEUTICALS 1008=2 55=19 1007=GBPCAD 1008=5 55=20499 1007=CENTURYLINK INC 1008=2 55=20500 1007=APPLE COMPUTER INC 1008=2 55=20501 1007=WR GRACE & CO 1008=2 55=22 1007=USDNOK 1008=4 55=20502 1007=GROUP 1 AUTOMOTIVE 1008=2 55=23 1007=AUDCHF 1008=5 55=20503 1007=TRW AUTOMOTIVE 1008=2 55=24 1007=USDMXN 1008=4 55=20504 1007=TIBCO SOFTWARE 1008=2 55=25 1007=GBPNZD 1008=4 55=20505 1007=GASTAR EXPLRTN 1008=2 55=20506 1007=MARRIOTT INTL 1008=2 55=27 1007=CADCHF 1008=5 55=20507 1007=BOSTON SCIEN CP 1008=2 55=20508 1007=GOODRICH PETRLM 1008=2 55=29 1007=USDSEK 1008=4 55=20509 1007=GOPRO 1008=2 55=20510 1007=AT&T 1008=2 55=20511 1007=TABLEAU SOFTWARE INC 1008=2 55=20512 1007=EPAM SYSTEMS INC 1008=2 55=20513 1007=D.R. HORTON INC 1008=2 55=20514 1007=BIO-REF LAB 1008=2 55=20515 1007=MCCORMICK & CO 1008=2 55=20516 1007=MEDTRONIC 1008=2 55=20517 1007=MONDELEZ INTL 1008=2 55=20518 1007=BIG 5 SPORTING GDS 1008=2 55=20519 1007=DEXCOM 1008=2 55=40 1007=GBPCHF 1008=5 55=20520 1007=TOTAL SYS SVCS 1008=2 55=41 1007=XAUUSD 1008=2 55=20521 1007=TUPPERWARE BRANDS CORPORATION 1008=2 55=42 1007=XAGUSD 1008=2 55=20522 1007=MEAD JOHNSON NUTRITION 1008=2 55=43 1007=KRWINR 1008=5 55=20523 1007=BLACKSTONE GROUP LP 1008=2 55=44 1007=EURRUB 1008=3 55=20524 1007=FRONTIER COMMUNICATIONS CORP 1008=2 55=45 1007=USDRUB 1008=3 55=20525 1007=MCKESSON CORP 1008=2 55=46 1007=USDCNH 1008=4 55=20526 1007=MELCO CROWN 1008=2 55=47 1007=EURSEK 1008=4 55=20527 1007=FLEETMATICS 1008=2 55=48 1007=GBPZAR 1008=4 55=20528 1007=DISCOVERY HOLDING CO 1008=2 55=49 1007=USDTRY 1008=4 55=20529 1007=EXELIXIS 1008=2 55=50 1007=NZDCHF 1008=4 55=20530 1007=DILLARD CL A 1008=2 55=51 1007=EURPLN 1008=4 55=20531 1007=LEVEL 3 COMMS 1008=2 55=52 1007=USDZAR 1008=5 55=20532 1007=FINISAR 1008=2 55=53 1007=EURMXN 1008=4 55=20533 1007=FIDELITY NATL IN 1008=2 55=54 1007=EURDKK 1008=4 55=20534 1007=ARRIS GROUP INC. 1008=2 55=55 1007=EURHUF 1008=2 55=20535 1007=HASBRO INC 1008=2 55=56 1007=SGDJPY 1008=3 55=20536 1007=GENERAL MILLS 1008=2 55=57 1007=EURNOK 1008=4 55=20537 1007=MARATHON PETROLEUM CORP 1008=2 55=58 1007=EURHKD 1008=4 55=20538 1007=MIDDLEBY CORP 1008=2 55=59 1007=GBPNOK 1008=4 55=20539 1007=BALL CORP 1008=2 55=60 1007=USDHUF 1008=3 55=20540 1007=MEREDITH CORP 1008=2 55=61 1007=NZDCAD 1008=4 55=20541 1007=DANGDANG INC 1008=2 55=62 1007=EURZAR 1008=5 55=20542 1007=ACCENTURE PLC 1008=2 55=63 1007=EURCZK 1008=4 55=20543 1007=AVANIR PHARMA 1008=2 55=64 1007=NZDSGD 1008=4 55=20544 1007=CENTENE GROUP 1008=2 55=65 1007=USDPLN 1008=4 55=20545 1007=CINEMARK HLDG 1008=2 55=66 1007=EURNZD 1008=4 55=20546 1007=ANI PHARMACEUTICALS 1008=2 55=67 1007=AUDNZD 1008=4 55=20547 1007=ALLEGION CLOSE ONLY 1008=2 55=68 1007=EURTRY 1008=4 55=20548 1007=GENERAL DYNAMICS 1008=2 55=69 1007=USDDKK 1008=4 55=20549 1007=MOODY 1008=2 55=70 1007=NZDJPY 1008=2 55=20550 1007=COVANCE INC 1008=2 55=71 1007=USDHKD 1008=4 55=20551 1007=CAN NATL RAILWAY 1008=2 55=72 1007=USDCZK 1008=4 55=20552 1007=CATALENT 1008=2 55=73 1007=USDSGD 1008=5 55=20553 1007=MGM RESORTS INTL 1008=2 55=74 1007=GBPSGD 1008=4 55=20554 1007=ADVANCED ENERGY 1008=2 55=20555 1007=DISH NETWORK CORP 1008=2 55=20556 1007=FLOWERS FOODS 1008=2 55=20557 1007=ANIKA THERAPEUTICS 1008=2 55=20558 1007=MFC INDUSTRIAL 1008=2 55=20559 1007=ACTAVIS INC 1008=2 55=20560 1007=CTC MEDIA INC 1008=2 55=20561 1007=CIGNA CORP 1008=2 55=20562 1007=ALTRIA GROUP 1008=2 55=20563 1007=CHICAGO BR & IRN 1008=2 55=20564 1007=ELDORADO GOLD 1008=2 55=20565 1007=T-MOBILE US 1008=2 55=20566 1007=DISCOVER FINANCL 1008=2 55=20567 1007=ADVANTAGE OIL GAS 1008=2 55=20568 1007=THE BUCKLE INC 1008=2 55=20569 1007=MASTERCARD CL A 1008=2 55=20570 1007=CB RICHARDS ELLIS 1008=2 55=20571 1007=ACACIA RESEARCH 1008=2 55=20572 1007=AGL RESOURCES INC 1008=2 55=20573 1007=MAZOR ROBOTICS 1008=2 55=20574 1007=STH WEST AIRLINES 1008=2 55=20575 1007=COGNIZANT TECHNOLOGY SOLUTIO 1008=2 55=20576 1007=MOBILEYE NV 1008=2 55=20577 1007=TAKE-TWO INTERACTIVE 1008=2 55=20578 1007=NABORS INDS 1008=2 55=20579 1007=MECOX LAND LTD 1008=2 55=20580 1007=NEWMONT MINING 1008=2 55=20581 1007=FREEPORT MCM 1008=2 55=20582 1007=ARES CAPITAL 1008=2 55=20583 1007=AVAGO TECHNOLOGIES 1008=2 55=20584 1007=THORATEC CORP 1008=2 55=20585 1007=EDWARDS LIFESCIENCES 1008=2 55=20586 1007=FLOWSERVE CORP 1008=2 55=20587 1007=MICROCHIP TECHNOLOGY INC 1008=2 55=20588 1007=MYLAN INC 1008=2 55=20589 1007=GENERAC HOLDINGS 1008=2 55=20590 1007=TETRALOGIC PHARMA 1008=2 55=20591 1007=TURKCELL ILETI 1008=2 55=20592 1007=GAP.INC 1008=2 55=20593 1007=CLEAN DIESEL TECH *CLOSE ONLY* 1008=2 55=20594 1007=CAN IMPL BK COMM 1008=2 55=20595 1007=MANNKIND CORP 1008=2 55=20596 1007=DU PONT (EI) 1008=2 55=20597 1007=FIRST AM FIN 1008=2 55=20598 1007=CITIGROUP 1008=2 55=20599 1007=MONSTER WORLDWIDE INC 1008=2 55=20600 1007=ALLIANCE DATA SYSTEMS 1008=2 55=20601 1007=FCB FINANCIAL HLDS 1008=2 55=20602 1007=EQUIFAX INC 1008=2 55=20603 1007=AMERISOURCEBERGN 1008=2 55=20604 1007=GULFPORT ENERGY 1008=2 55=20605 1007=ENCANA CORP 1008=2 55=20606 1007=NASDAQ OMX GROUP 1008=2 55=20607 1007=TJX CO INC 1008=2 55=20608 1007=CYTEC INDS 1008=2 55=20609 1007=ENTEROMEDICS INC 1008=2 55=20610 1007=ENTERPRISE PRODUCTS 1008=2 55=20611 1007=TWENTY-FIRST CEN. FOX 1008=2 55=20612 1007=NICE SYSTEMS LTD 1008=2 55=20613 1007=GILEAD SCIENCES INC 1008=2 55=20614 1007=CHICAGO MERCANTL 1008=2 55=20615 1007=COBALT INT. ENERGY 1008=2 55=20616 1007=CREDIT SUISSE 1008=2 55=20617 1007=GOOGLE CLASS C 1008=2 55=20618 1007=ALLERGAN INC 1008=2 55=20619 1007=TROVAGENE INC 1008=2 55=20620 1007=BIG LOTS INC 1008=2 55=20621 1007=MARVELL TECHNOLOGY GROUP LTD 1008=2 55=20622 1007=GOGO 1008=2 55=20623 1007=CORELOGIC INC 1008=2 55=20624 1007=GENTHERM INC 1008=2 55=20625 1007=CHIMERA INVESTMENT CORP 1008=2 55=20626 1007=MOTOROLA SOLUTIONS 1008=2 55=20627 1007=ASPEN TECH 1008=2 55=20628 1007=LEUCADIA NATL 1008=2 55=20629 1007=BRISTOL-MYER SQUIIB 1008=2 55=20630 1007=DEAN FOODS CO 1008=2 55=20631 1007=MICROSOFT CORP 1008=2 55=20632 1007=NRG ENERGY INC 1008=2 55=20633 1007=NEWFIELD EXPLORATION CO 1008=2 55=20634 1007=METHANEX 1008=2 55=20635 1007=AARON 1008=2 55=20636 1007=A.O SMITH CORP 1008=2 55=20637 1007=ABERC FITCH A 1008=2 55=20638 1007=ASBURY AUTOMOTIVE 1008=2 55=20639 1007=TEXTAINER GROUP 1008=2 55=20640 1007=BRIGHTCOVE INC 1008=2 55=20641 1007=DECKERS OUTDOOR CORP 1008=2 55=20642 1007=APOLLO GROUP INC 1008=2 55=20643 1007=TRAVIS PERKINS 1008=2 55=20644 1007=CONTINENTAL AG 1008=2 55=20645 1007=KESKO 1008=2 55=20646 1007=LADBROKES 1008=2 55=20647 1007=DERWENT LONDON 1008=2 55=20648 1007=PAYPOINT 1008=2 55=20649 1007=ARSEUS NV 1008=2 55=20650 1007=SKF B 1008=2 55=20651 1007=DUNELM GROUP PLC 1008=2 55=20652 1007=PEUGEOT 1008=3 55=20653 1007=TELEFONICA DE 1008=2 55=20654 1007=BOUYGUES 1008=3 55=20655 1007=ASHTEAD GROUP 1008=2 55=20656 1007=MERCIALYS 1008=2 55=20657 1007=TULLOW OIL 1008=2 55=20658 1007=EDP 1008=2 55=20659 1007=ASOS PLC 1008=2 55=20660 1007=SEMAPA 1008=2 55=20661 1007=GUERBET 1008=2 55=20662 1007=NMC HEALTH PLC 1008=2 55=20663 1007=DE LA RUE 1008=2 55=20664 1007=HIGHCO 1008=2 55=20665 1007=BELLWAY 1008=2 55=20666 1007=STEF 1008=2 55=20667 1007=VEOLIA ENVIRON 1008=2 55=20668 1007=PRAKTIKER CLOSE ONLY 1008=2 55=20669 1007=TNT EXPRESS NV 1008=2 55=20670 1007=LATECOERE WRT (CLOSING ONLY) 1008=2 55=20671 1007=HEIDELBERGCEMENT 1008=2 55=20672 1007=TRIGANO 1008=2 55=20673 1007=SAMPO PLC 1008=2 55=20674 1007=BAYER AG 1008=2 55=20675 1007=WILLIAM HILL 1008=2 55=20676 1007=DERICHEBOURG 1008=2 55=20677 1007=RIGHTMOVE 1008=2 55=20678 1007=LAND SECURITIES 1008=2 55=20679 1007=RUBIS 1008=2 55=20680 1007=AXEL SPRINGER AG 1008=2 55=20681 1007=FRESNILLO PLC 1008=2 55=20682 1007=STALLARGENES 1008=2 55=20683 1007=INMARSAT 1008=2 55=20684 1007=OPHIR ENERGY 1008=2 55=20685 1007=DEVOTEAM 1008=2 55=20686 1007=PRUDENTIAL 1008=2 55=20687 1007=SVG CAPITAL ORD ?1 1008=2 55=20688 1007=SAINSBURY 1008=2 55=20689 1007=NEXT 1008=2 55=20690 1007=WOLTERS KLUWER 1008=2 55=20691 1007=L.V.M.H. 1008=2 55=20692 1007=SMITHS GROUP 1008=2 55=20693 1007=WHITBREAD 1008=2 55=20694 1007=SBM OFFSHORE 1008=3 55=20695 1007=PLUS500 1008=2 55=20696 1007=DAMARTEX 1008=2 55=20697 1007=SCOTTISH MORTGAGE INV TRUST 1008=2 55=20698 1007=ASSYSTEM 1008=2 55=20699 1007=WENDEL 1008=2 55=20700 1007=DEXIA 1008=2 55=20701 1007=TELECITY 1008=2 55=20702 1007=SDC - INVESTIMENTOS 1008=2 55=20703 1007=INVERKO 1008=2 55=20704 1007=BLACKROCK WORLD MINING TRUST PLC 1008=2 55=20705 1007=RENISHAW 1008=2 55=20706 1007=DIGNITY 1008=2 55=20707 1007=WPP 1008=2 55=20708 1007=BONDUELLE 1008=2 55=20709 1007=WINCOR NIXDORF 1008=2 55=20710 1007=UNILEVER DR 1008=2 55=20711 1007=GALP ENERGIA 1008=2 55=20712 1007=SMITH AND NEPHEW 1008=2 55=20713 1007=BOIRON 1008=2 55=20714 1007=GENEL
ENERGY 1008=2 55=20715 1007=VICAT 1008=2 55=20716 1007=UNILEVER 1008=2 55=20717 1007=AFFINE R.E. 1008=2 55=20718 1007=ORANGE 1008=3 55=20719 1007=SWORD GROUP 1008=2 55=20720 1007=ZODIAC 1008=2 55=20721 1007=SWEDISH MATCH 1008=2 55=20722 1007=GENERAL SANTE (NO SHORTING) 1008=2 55=20723 1007=TSB BANK 1008=2 55=20724 1007=ENQUEST PLC 1008=2 55=20725 1007=CITY OF LONDON INVESTMENT TRUST 1008=2 55=20726 1007=R.E.N 1008=2 55=20727 1007=HARGREAVES 1008=2 55=20728 1007=SCHRODERS 1008=2 55=20729 1007=HOLLAND COL 1008=2 55=20730 1007=SAMSE 1008=2 55=20731 1007=ANTOFAGASTA 1008=2 55=20732 1007=UK COMM PROP TRUST 1008=2 55=20733 1007=NEOPOST 1008=2 55=20734 1007=POLAR CAPITAL TECHNOLOGY TRUST PLC 1008=2 55=20735 1007=WACKER CHEMIE 1008=2 55=20736 1007=RHEINMETALL 1008=3 55=20737 1007=JOHN LAING INFRA FUND 1008=2 55=20738 1007=MEDIA 6 1008=2 55=20739 1007=FINANCIERE ODET 1008=2 55=20740 1007=HENNES & MAURITZ 1008=2 55=20741 1007=PFEIFFER VACUUM 1008=2 55=20742 1007=WEIR GROUP 1008=2 55=20743 1007=OXFORD INSTRUMENTS PLC 1008=2 55=20744 1007=HEIDELBERG 1008=2 55=20745 1007=TFF GROUP 1008=2 55=20746 1007=SAFRAN 1008=3 55=20747 1007=CORIO 1008=2 55=20748 1007=HAVAS 1008=2 55=20749 1007=HANNOVER RUECK 1008=2 55=20750 1007=TELEPERFORMANCE 1008=2 55=20751 1007=ESSO 1008=2 55=20752 1007=LISI 1008=2 55=20753 1007=REGUS GROUP 1008=2 55=20754 1007=GREENCORE GROUP PLC 1008=2 55=20755 1007=GECINA 1008=2 55=20756 1007=UNITED BUSINESS MEDIA 1008=2 55=20757 1007=DELHAIZE 1008=3 55=20758 1007=BNP PARIBAS 1008=3 55=20759 1007=BBA GROUP 1008=2 55=20760 1007=KORIAN CLOSE ONLY 1008=2 55=20761 1007=SNOWWORLD 1008=2 55=20762 1007=FIDELITY EUROPEAN 1008=2 55=20763 1007=INTL PRSNL FIN 1008=2 55=20764 1007=MANITOU 1008=2 55=20765 1007=ABCAM PLC 1008=2 55=20766 1007=STANDARD LIFE 1008=2 55=20767 1007=AXA 1008=3 55=20768 1007=EDENRED 1008=2 55=20769 1007=PENNON 1008=2 55=20770 1007=NOKIA 1008=2 55=20771 1007=BRUNEL 1008=2 55=20772 1007=IMI 1008=2 55=20773 1007=HAMMERSON 1008=2 55=20774 1007=SIG 1008=2 55=20775 1007=REXEL 1008=3 55=20776 1007=ING GROEP 1008=2 55=20777 1007=EDINBURGH INVESTMENT TRUST 1008=2 55=20778 1007=MARSTONS 1008=2 55=20779 1007=INTU PROPERTIES 1008=2 55=20780 1007=BWIN.PARTY DIGITAL 1008=2 55=20781 1007=SALZGITTER 1008=2 55=20782 1007=BRITISH EMPIRE TRUST 1008=2 55=20783 1007=SEQUANA CAPITAL 1008=2 55=20784 1007=REED ELSEVIER 1008=2 55=20785 1007=NORTHGATE 1008=2 55=20786 1007=ALCATEL 1008=2 55=20787 1007=EUROFINS CEREP 1008=2 55=20788 1007=TR PROPERTY INVESTMENT TST 1008=2 55=20789 1007=PERFORM GROUP 1008=2 55=20790 1007=DANONE 1008=3 55=20791 1007=SAVILLS 1008=2 55=20792 1007=TESCO 1008=2 55=20793 1007=BUNZL 1008=2 55=20794 1007=BOSKALIS WESTMIN 1008=2 55=20795 1007=DAIMLER AG 1008=2 55=20796 1007=SIEMENS 1008=3 55=20797 1007=FIDESSA GROUP PLC 1008=2 55=20798 1007=ASM INTERNATIONAL 1008=2 55=20799 1007=MERCHANTS TRUST 1008=2 55=20800 1007=RSA INSURANCE GROUP 1008=2 55=20801 1007=HERMES INTERNATIONAL 1008=2 55=20802 1007=FIMALAC 1008=2 55=20803 1007=DEUTSCHE POST 1008=2 55=20804 1007=ABB 1008=2 55=20805 1007=COFINIMMO 1008=3 55=20806 1007=SAGE GROUP 1008=2 55=20807 1007=ALTEN 1008=3 55=20808 1007=CMB 1008=2 55=20809 1007=ETAM DEVELOP 1008=2 55=20810 1007=PHOENIX GROUP 1008=2 55=20811 1007=ARCADIS 1008=3 55=20812 1007=EXPERIAN GRP 1008=2 55=20813 1007=PZ CUSSONS 1008=2 55=20814 1007=HAMBURGER HAFEN 1008=2 55=20815 1007=SOPRA GROUP 1008=2 55=20816 1007=HAYS 1008=2 55=20817 1007=TELE2 1008=2 55=20818 1007=ROBERTET 1008=2 55=20819 1007=CREST NICHOLSON 1008=2 55=20820 1007=AKZO NOBEL 1008=2 55=20821 1007=INGENICO 1008=2 55=20822 1007=BRILL 1008=2 55=20823 1007=LECTRA 1008=2 55=20824 1007=SCHNEIDER ELECTR 1008=2 55=20825 1007=MANULTAN 1008=2 55=20826 1007=SAP AG 1008=2 55=20827 1007=ROTORK 1008=2 55=20828 1007=JOHNSON MATTHEY 1008=2 55=20829 1007=CASINO 1008=3 55=20830 1007=AVENIR TELECOM 1008=2 55=20831 1007=WITAN INVESTMENT COMPANY 1008=2 55=20832 1007=KELLER GROUP 1008=2 55=20833 1007=AVIVA 1008=2 55=20834 1007=TUI TRAVEL 1008=2 55=20835 1007=DUERR AG 1008=2 55=20836 1007=ZALANDO 1008=2 55=20837 1007=TAYLOR WIMPEY 1008=2 55=20838 1007=HISCOX 1008=2 55=20839 1007=EXACT 1008=2 55=20840 1007=HERALD INV TR COMMON STOCK 1008=2 55=20841 1007=UPONOR 1008=2 55=20842 1007=CARGOTEC CORP 1008=2 55=20843 1007=LOCINDUS 1008=2 55=20844 1007=DEVRO PLC COMMON STOCK 1008=2 55=20845 1007=MPI 1008=2 55=20846 1007=BENETEAU 1008=3 55=20847 1007=HUNTING 1008=2 55=20848 1007=BSKYB 1008=2 55=20849 1007=TALVIVAARAN MINING 1008=2 55=20850 1007=BRICORAMA 1008=2 55=20851 1007=BOLIDEN 1008=2 55=20852 1007=CSR 1008=2 55=20853 1007=HOMESERVE 1008=2 55=20854 1007=IMAGINATION TECHNOLOGIES GROUP PLC 1008=2 55=20855 1007=INVESTOR 1008=2 55=20856 1007=STAGECOACH GROUP 1008=2 55=20857 1007=MILLENNIUM & COPTHORNE 1008=2 55=20858 1007=LAFUMA 1008=2 55=20859 1007=ROYAL DUTCH SHELL A SHR 1008=2 55=20860 1007=AVEVA GROUP 1008=2 55=20861 1007=MONKS INVESTMENT TRUST 1008=2 55=20862 1007=ASML HOLDING 1008=2 55=20863 1007=GREAT PORTLAND ESTATES 1008=2 55=20864 1007=EDF 1008=3 55=20865 1007=HEIJMANS 1008=2 55=20866 1007=LANXESS 1008=3 55=20867 1007=BRITVIC 1008=2 55=20868 1007=VAN LANSCHOT 1008=2 55=20869 1007=ROLLS-ROYCE 1008=2 55=20870 1007=BRITISH AMERICAN TOBACCO 1008=2 55=20871 1007=ALBIOMA 1008=2 55=20872 1007=ZIGGO 1008=2 55=20874 1007=BANCO BPI 1008=2 55=20875 1007=PREMIER OIL 1008=2 55=20876 1007=BODYCOTE 1008=2 55=20877 1007=PISC DESJOYAUX 1008=2 55=20878 1007=EURO RESSOURCES 1008=2 55=20879 1007=STRATEC BIOMEDICAL 1008=2 55=20880 1007=RANDSTAD 1008=3 55=20881 1007=TETRAGON FINANCIAL 1008=2 55=20882 1007=TELEGRAAF 1008=2 55=20883 1007=AVANQUEST SOFTWARE 1008=2 55=20884 1007=QIAGEN 1008=2 55=20885 1007=SABMILLER 1008=2 55=20886 1007=MORGAN ADVANCED 1008=2 55=20887 1007=HENDERSON GROUP 1008=2 55=20888 1007=GKN 1008=2 55=20889 1007=CELESIO 1008=3 55=20890 1007=FUGRO 1008=2 55=20891 1007=INTL PUBLIC PARTNERSHIP LTD 1008=2 55=20892 1007=BARCO 1008=2 55=20893 1007=QINETIQ 1008=2 55=20894 1007=SONAE INDUSTRIA 1008=2 55=20895 1007=RHOEN KLINIKUM 1008=2 55=20896 1007=SKANSKA B 1008=2 55=20897 1007=NESTE OIL 1008=2 55=20898 1007=NEXTRADIOTV 1008=2 55=20899 1007=NICOX 1008=2 55=20900 1007=ICADE 1008=3 55=20901 1007=SAGA GRP 1008=2 55=20902 1007=KONECRANES 1008=2 55=20903 1007=LANSON-BCC 1008=2 55=20905 1007=AURUBIS AG 1008=2 55=20906 1007=LONDONMETRIC 1008=2 55=20907 1007=GREENE KING 1008=2 55=20908 1007=PETRA DIAMONDS LIMITED 1008=2 55=2]

*/