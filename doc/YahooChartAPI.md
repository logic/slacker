Originally from:

http://code.google.com/p/yahoo-finance-managed/wiki/miscapiImageDownload

Mirrored here, since `code.google.com` is shutting down.

----

# Chart Download

How to download technical analysing charts of a stock, index or currency
exchange.

## Start

http://chart.finance.yahoo.com/z?

## ID

Now, you have to set the IDs you want to receive. Every stock, index or
currency has their own ID.

You also have to convert special characters into the correct URL format

Add to the end `s=` and the ID.

http://chart.finance.yahoo.com/z?s=GOOG

## Time span

After that you have to assign the size of the investigated period.

Add to the end `&t=` and the time span tag.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m

## Type

Now, add to the end `&q=` and the chart type tag.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l

## Scale

For defining the scaling type just add `&l=` and the chart scale tag.

Here you only have to describe, if logarithmic scaling is `on` or `off`.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on

## Size

Add to the end `&z=` and the chart size tag of middle or large.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on&z=l

## Moving Average Indicator

Here you can declare the moving average intervals.

If you have more than one seperate them with `,`.

Add to the end `&p=` and for each interval `m` and the interval tag.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on&z=l&p=m50,m200

## Exponential Moving Average Indicator

Here you can declare the exponential moving average intervals. It's nearly the
same indicator like moving average, but its values are calculated with the
standard deviation.

If you have more than one seperate them with `,`.

If you didn't declare a moving average interval one step before, add to the end
`&p=` and for each interval `e` and the interval tag.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on&z=l&p=e50,e200

If you did, just add the comma seperated tags.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on&z=l&p=m50,e200

## Technical Indicators 1

Here you can declare the first part of technical indicators.

If you have more than one seperate them with `,`.

It's the same rule like with exponential moving average.

This part is valid for the following indicators:

* Bollinger_Bands
* Parabolic_SAR
* Splits
* Volume

If you didn't declare a moving or exponential moving average interval one or
two steps before, add to the end `&p=` and for each indicator the indicator
tag.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on&z=l&p=v,b

If you did, just add the comma seperated tags.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on&z=l&p=m50,e200,v

## Technical Indicators 2

Here you can declare the second part of technical indicators.

If you have more than one seperate them with `,`.

This part is valid for the following indicators:

* MACD
* MFI
* ROC
* RSI
* Slow_Stoch
* Fast_Stoch
* Vol_MA
* W_R

Add to the end of the URL `&a=` and for each indicator the tag.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on&z=l&p=m50,e200,v&a=p12,p

## Comparing IDs

Lastly you are able to declare the the IDs of different stocks or indices that
will be compared with the base ID of this image.

If you compare the stocks or indices, the image will display the chart scaling
in relative values in percent, not absolute in any currency unit.

Comparing currency exchanges is also possible.

If you have more than one seperate them with `,`.

You also have to convert special characters into the correct URL format

Add to the end of the URL "&c=" and for each stock, index or currency the ID.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on&z=l&p=m50,e200,v&a=p12,p&c=%5EDJI,EURUSD=X

## Culture

Add to the end &region= and the culture tag of a country of your choice.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on&z=l&p=m50,e200,v&a=p12,p&c=%5EDJI,EURUSD=X&region=DE

Also add to the end of URL &lang= and a language/country of your choice.

http://chart.finance.yahoo.com/z?s=GOOG&t=6m&q=l&l=on&z=l&p=m50,e200,v&a=p12,p&c=%5EDJI,EURUSD=X&region=DE&lang=de-DE
