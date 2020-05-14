package backtest

import (
	. "github.com/coinrust/crex"
	"time"
)

const reportHistoryTemplate = `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=UTF-8">
    <title>Trade History Report</title>
    <meta name="generator" content="client terminal">
    <style type="text/css">
    <!--
    @media screen {
      td { font: 8pt  Tahoma,Arial; }
      th { font: 10pt Tahoma,Arial; }
    }
    @media print {
      td { font: 7pt Tahoma,Arial; }
      th { font: 9pt Tahoma,Arial; }
    }
    .msdate { mso-number-format:"General Date"; }
    .mspt   { mso-number-format:\#\,\#\#0\.00;  }
    body {margin:1px;}
    //-->
    </style>
  </head>
<body>
<div align="center">
<table cellspacing="1" cellpadding="3" border="0">
   <tr align="center">
      <td colspan="13"><div style="font: 14pt Tahoma"><b>Trade History Report</b><br></div></td>
   </tr>
   <!--<tr align="left">
       <th colspan="3" nowrap align="right" style="width: 220px; height: 20px">Date:</th>
       <th colspan="10" nowrap align="left" style="width: 220px; height: 20px"><b>2020.05.14 11:15</b></th>
   </tr>-->
   <tr>
      <td nowrap style="width: 140px;height: 10px"></td>
      <td nowrap style="width: 60px;"></td>
      <td nowrap style="width: 60px;"></td>
      <td nowrap style="width: 60px;"></td>
      <td nowrap style="width: 60px;"></td>
      <td nowrap style="width: 60px;"></td>
      <td nowrap style="width: 60px;"></td>
      <td nowrap style="width: 70px;"></td>
      <td nowrap style="width: 70px;"></td>
      <td nowrap style="width: 70px;"></td>
      <td nowrap style="width: 70px;"></td>
      <td nowrap style="width: 60px;"></td>
      <td nowrap style="width: 100px;"></td>
   </tr>
   <tr align="center">
      <th colspan="13" style="height: 25px"><div style="font: 10pt Tahoma"><b>Orders</b></div></th>
   </tr>
   <tr align="center" bgcolor="#E5F0FC">
      <td nowrap style="height: 30px"><b>Open Time</b></td>
      <td nowrap><b>Order</b></td>
      <td nowrap><b>Symbol</b></td>
      <td nowrap><b>Type</b></td>
      <td nowrap colspan="2"><b>Volume</b></td>
      <td nowrap><b>Price</b></td>
      <td nowrap><b>S / L</b></td>
      <td nowrap><b>T / P</b></td>
      <td nowrap colspan="2"><b>Time</b></td>
      <td nowrap><b>State</b></td>
      <td nowrap><b>Comment</b></td>
   </tr>
   <!--{order-row}-->
   <tr>
      <td nowrap style="height: 10px"></td>
   </tr>
   <tr align="center">
      <th colspan="13" style="height: 25px"><div style="font: 10pt Tahoma"><b>Deals</b></div></th>
   </tr>
   <tr align="center" bgcolor="#E5F0FC">
      <td nowrap style="height: 30px"><b>Time</b></td>
      <td nowrap><b>Deal</b></td>
      <td nowrap><b>Symbol</b></td>
      <td nowrap><b>Type</b></td>
      <td nowrap><b>Direction</b></td>
      <td nowrap><b>Volume</b></td>
      <td nowrap><b>Price</b></td>
      <td nowrap><b>Order</b></td>
      <td nowrap><b>Commission</b></td>
      <td nowrap><b>Swap</b></td>
      <td nowrap><b>Profit</b></td>
      <td nowrap><b>Balance</b></td>
      <td nowrap><b>Comment</b></td>
   </tr>
   <!--
   <tr bgcolor="#FFFFFF" align="right"><td nowrap>2018.06.01 03:08:53</td><td nowrap>10320598</td><td nowrap></td><td nowrap>balance</td><td nowrap></td><td nowrap></td><td nowrap></td><td nowrap></td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>50 000.00</td><td nowrap>50 000.00</td><td nowrap></td></tr>
   <tr bgcolor="#F7F7F7" align="right"><td nowrap>2018.12.03 09:45:34</td><td nowrap>10890657</td><td nowrap>USDCHF</td><td nowrap>sell</td><td nowrap>in</td><td nowrap>0.20</td><td nowrap>0.99669</td><td nowrap>11984150</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>50 000.00</td><td nowrap></td></tr>
   <tr bgcolor="#FFFFFF" align="right"><td nowrap>2018.12.03 09:49:01</td><td nowrap>10890672</td><td nowrap>USDCHF</td><td nowrap>buy</td><td nowrap>out</td><td nowrap>0.20</td><td nowrap>0.99667</td><td nowrap>11984165</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>0.40</td><td nowrap>50 000.40</td><td nowrap></td></tr>
   <tr bgcolor="#F7F7F7" align="right"><td nowrap>2018.12.03 09:50:59</td><td nowrap>10890673</td><td nowrap>USDCHF</td><td nowrap>buy</td><td nowrap>in</td><td nowrap>0.20</td><td nowrap>0.99652</td><td nowrap>11984167</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>50 000.40</td><td nowrap></td></tr>
   <tr bgcolor="#FFFFFF" align="right"><td nowrap>2018.12.03 09:51:34</td><td nowrap>10890676</td><td nowrap>USDCHF</td><td nowrap>sell</td><td nowrap>out</td><td nowrap>0.20</td><td nowrap>0.99661</td><td nowrap>11984171</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>1.81</td><td nowrap>50 002.21</td><td nowrap></td></tr>
   <tr bgcolor="#F7F7F7" align="right"><td nowrap>2018.12.03 09:53:44</td><td nowrap>10890681</td><td nowrap>USDCHF</td><td nowrap>sell</td><td nowrap>in</td><td nowrap>0.20</td><td nowrap>0.99656</td><td nowrap>11984176</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>50 002.21</td><td nowrap></td></tr>
   <tr bgcolor="#FFFFFF" align="right"><td nowrap>2018.12.03 09:54:12</td><td nowrap>10890683</td><td nowrap>USDCHF</td><td nowrap>buy</td><td nowrap>out</td><td nowrap>0.20</td><td nowrap>0.99655</td><td nowrap>11984178</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>0.20</td><td nowrap>50 002.41</td><td nowrap></td></tr>
   <tr bgcolor="#F7F7F7" align="right"><td nowrap>2018.12.03 09:56:49</td><td nowrap>10890689</td><td nowrap>USDCHF</td><td nowrap>buy</td><td nowrap>in</td><td nowrap>0.20</td><td nowrap>0.99658</td><td nowrap>11984187</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>50 002.41</td><td nowrap></td></tr>
   <tr bgcolor="#FFFFFF" align="right"><td nowrap>2018.12.03 09:57:07</td><td nowrap>10890693</td><td nowrap>USDCHF</td><td nowrap>sell</td><td nowrap>out</td><td nowrap>0.20</td><td nowrap>0.99663</td><td nowrap>11984191</td><td nowrap>0.00</td><td nowrap>0.00</td><td nowrap>1.00</td><td nowrap>50 003.41</td><td nowrap></td></tr>
   -->
   <!--
   <tr align="right">
      <td nowrap colspan="8" style="height: 30px"></td>
      <td nowrap><b>0.00</b></td>
      <td nowrap><b>0.00</b></td>
      <td nowrap><b>3.41</b></td>
      <td nowrap><b>50 003.41</b></td>
      <td nowrap></td>
   </tr>
   -->
   <tr align="right">
      <td colspan="13" style="height: 10px"></td>
   </tr>
   <tr align="right">
      <td colspan="13" style="height: 10px"></td>
   </tr>
</table>
</div>
</body>
</html>`

// SOrder "event":"order"/"deal"
type SOrder struct {
	Ts        time.Time  // ts: 2019-10-02T07:03:53.584+0800
	Order     *Order     // order
	OrderBook *OrderBook // orderbook
	Comment   string     // msg: Place order/Match order
}
