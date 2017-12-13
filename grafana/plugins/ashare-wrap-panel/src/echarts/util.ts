const upColor: string = '#ec0000';
const upBorderColor: string = '#8A0000';
const downColor: string = '#00da3c';
const downBorderColor: string = '#008F28';

export class EchartsUtil {
  // 数据意义：开盘(open),收盘(close),最低(low),最高(high),成交量(volume),dea/dif/macd(macd)
  private rawShareData: any[] = [];

  constructor() {
  }

  setData(rawShareData: any[]) {
    this.rawShareData = rawShareData
  }

  splitShareData() {
    var categoryData = new Array();
    var values = new Array();
    var volumes = new Array();
    var macdDIFs = new Array();
    var macdDEAs = new Array();
    var macdPillers = new Array();
    var data = this.rawShareData;
    for (var i = 0; i < data.length; i++) {
      categoryData.push(data[i].splice(0, 1)[0]);
      values.push(data[i])
      volumes.push([i, data[i][4], data[i][0] > data[i][1] ? 1 : -1]);
      macdDIFs.push(data[i][5]);
      macdDEAs.push(data[i][6]);
      macdPillers.push([i, data[i][7], data[i][7] >= 0 ? 1 : -1]);
    }
    return {
        categoryData: categoryData,
        values: values,
        volumes: volumes,
        macdDIFs: macdDIFs, 
        macdDEAs: macdDEAs, 
        macdPillers: macdPillers, 
    };
  }

  calculateMA(dayCount: number, splitData): any[] {
    var result = new Array();
    for (var i = 0, len = splitData.values.length; i < len; i++) {
        if (i < dayCount) {
            result.push('-');
            continue;
        }
        var sum = 0;
        for (var j = 0; j < dayCount; j++) {
          sum += parseInt(splitData.values[i - j][1]);
        }
        result.push((sum / dayCount).toFixed(2));
    }
    return result;
  }

  getKLineOption() {
    var splitData = this.splitShareData();
    return {
      backgroundColor: '#fff',
      title: {
          text: '上证指数',
          left: 0
      },
      animation: false,
      legend: {
          // bottom: 10,
          left: 'center',
          data: ['日K', 'MA5', 'MA10', 'MA20', 'MA30']
      },
      tooltip: {
          trigger: 'axis',
          axisPointer: {
              type: 'cross'
          },
          backgroundColor: 'rgba(245, 245, 245, 0.8)',
          borderWidth: 1,
          borderColor: '#ccc',
          padding: 10,
          textStyle: {
              color: '#000'
          },
          position: function (pos, params, el, elRect, size) {
              var obj = {top: 10};
              obj[['left', 'right'][+(pos[0] < size.viewSize[0] / 2)]] = 30;
              return obj;
          }
      },
      axisPointer: {
        link: {xAxisIndex: 'all'},
        label: {
          backgroundColor: '#777'
        }
      },
      grid: [
        {
          left: '10%',
          right: '8%',
          height: '40%'
        },
        {
          left: '10%',
          right: '8%',
          top: '55%',
          height: '20%'
        },
        {
          left: '10%',
          right: '8%',
          top: '75%',
          height: '15%'
        },
      ],
      xAxis: [
        {
          type: 'category',
          data: splitData.categoryData,
          scale: true,
          boundaryGap : false,
          axisLine: {onZero: false},
          splitLine: {show: false},
          splitNumber: 20,
          min: 'dataMin',
          max: 'dataMax'
        },
        {
          type: 'category',
          data: splitData.categoryData,
          gridIndex: 1,
          scale: true,
          boundaryGap : false,
					silent: false,
          axisTick: {show: false},
					splitLine: {show: false},
          axisLabel: {show: false},
        },
        {
          type: 'category',
          gridIndex: 2,
          data: splitData.categoryData,
          scale: true,
          boundaryGap : false,
          axisLine: {onZero: false},
          axisTick: {show: false},
          splitLine: {show: false},
          axisLabel: {show: false},
          splitNumber: 20,
          min: 'dataMin',
          max: 'dataMax'
        },
      ],
      yAxis: [
	  		{
          scale: true,
          splitArea: {
            show: true
          }
        },
        {
          gridIndex: 1,
          scale: true,
          splitNumber: 2,
          axisLabel: {show: false},
          axisLine: {show: false},
          axisTick: {show: false},
          splitLine: {show: false}
        },
        {
          gridIndex: 2,
          scale: true,
          splitNumber: 2,
          axisLabel: {show: false},
          axisLine: {show: false},
          axisTick: {show: false},
          splitLine: {show: false}
        }
	  	],
      dataZoom: [
        {
          type: 'inside',
          xAxisIndex: [0, 1, 2],
          start: 90,
          end: 100
        },
        {
          show: true,
          type: 'slider',
          xAxisIndex: [0, 1, 2],
          top: '90%',
          start: 90,
          end: 100
        }
      ],
      series: [
        {
          name: '日K',
          type: 'candlestick',
          data: splitData.values,
          itemStyle: {
            normal: {
              color: upColor,
              color0: downColor,
              borderColor: upBorderColor,
              borderColor0: downBorderColor
            }
          },
          // markPoint: {
          //     label: {
          //         normal: {
          //             formatter: function (param) {
          //                 return param != null ? Math.round(param.value) : '';
          //             }
          //         }
          //     },
          //     data: [
          //         {
          //             name: 'XX标点',
          //             coord: ['2013/5/31', 2300],
          //             value: 2300,
          //             itemStyle: {
          //                 normal: {color: 'rgb(41,60,85)'}
          //             }
          //         },
          //         {
          //             name: 'highest value',
          //             type: 'max',
          //             valueDim: 'highest'
          //         },
          //         {
          //             name: 'lowest value',
          //             type: 'min',
          //             valueDim: 'lowest'
          //         },
          //         {
          //             name: 'average value on close',
          //             type: 'average',
          //             valueDim: 'close'
          //         }
          //     ],
          //     tooltip: {
          //         formatter: function (param) {
          //             return param.name + '<br>' + (param.data.coord || '');
          //         }
          //     }
          // },
          // markLine: {
          //     symbol: ['none', 'none'],
          //     data: [
          //         [
          //             {
          //                 name: 'from lowest to highest',
          //                 type: 'min',
          //                 valueDim: 'lowest',
          //                 symbol: 'circle',
          //                 symbolSize: 10,
          //                 label: {
          //                     normal: {show: false},
          //                     emphasis: {show: false}
          //                 }
          //             },
          //             {
          //                 type: 'max',
          //                 valueDim: 'highest',
          //                 symbol: 'circle',
          //                 symbolSize: 10,
          //                 label: {
          //                     normal: {show: false},
          //                     emphasis: {show: false}
          //                 }
          //             }
          //         ],
          //         {
          //             name: 'min line on close',
          //             type: 'min',
          //             valueDim: 'close'
          //         },
          //         {
          //             name: 'max line on close',
          //             type: 'max',
          //             valueDim: 'close'
          //         }
          //     ]
          // }
        },
        {
          name: 'MA5',
          type: 'line',
          data: this.calculateMA(5, splitData),
          smooth: true,
          lineStyle: {
            normal: {opacity: 0.5}
          }
        },
        {
          name: 'MA10',
          type: 'line',
          data: this.calculateMA(10, splitData),
          smooth: true,
          lineStyle: {
            normal: {opacity: 0.5}
          }
        },
        {
          name: 'MA10',
          type: 'line',
          data: this.calculateMA(10, splitData),
          smooth: true,
          lineStyle: {
            normal: {opacity: 0.5}
          }
        },
        {
          name: 'MACD',
          type: 'bar',
          xAxisIndex: 1,
          yAxisIndex: 1,
          data: splitData.macdPillers,
        },
        {
          name: 'VOL',
          type: 'bar',
          xAxisIndex: 2,
          yAxisIndex: 2,
          data: splitData.volumes
        },
        {
          name: 'DIF',
          type: 'line',
          xAxisIndex: 1,
          yAxisIndex: 1,
          data: splitData.macdDIFs,
					itemStyle: {
							normal: {
									color: 'rgb(0, 0, 0)'
							}
					},
          smooth: true,
          lineStyle: {
            normal: {opacity: 0.5}
          }
        },
        {
          name: 'DEA',
          type: 'line',
          xAxisIndex: 1,
          yAxisIndex: 1,
          data: splitData.macdDEAs,
					itemStyle: {
							normal: {
									color: 'rgb(138,43,226)'
							}
					},
          smooth: true,
        },
      ],
      visualMap: {
        show: false,
        seriesIndex: [3, 4],
        dimension: 2,
        pieces: [
          {
            value: 1,
            color: downColor
          },
          {
            value: -1,
            color: upColor
          }
        ]
      },
    }
  }
}
