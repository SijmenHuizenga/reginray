import React, {Component} from 'react';
import { ContextMenu, Item, ContextMenuProvider } from 'react-contexify';
import 'react-contexify/dist/ReactContexify.min.css';
import {BarChart, ChartContainer, ChartRow, Charts, Resizable, YAxis,} from "react-timeseries-charts";
import {TimeRange, TimeSeries} from "pondjs";


const upDownStyle = {
  value : {
    normal: {fill: "white", opacity: 1},
    highlighted: {fill: "white"},
    selected: {fill: "white"},
    muted: {fill: "white", opacity: 0.0},
  }
};

export class TopBar extends Component {

  render() {

    let a = this.makeTimeseries();
    let series = new TimeSeries({
      name : `Log Series`,
      columns : ["index", "value"],
      points: a
    });

    return <div>
      <ContextMenuProvider id="timeselectormenu">
        <TimeSelector trafficSeries={series} timeRangeChanged={this.props.timeRangeChanged}/>
      </ContextMenuProvider>
      {this.myAwesomeMenu()}
    </div>
  }

  onClick({ event, ref, data }){
    let now = unixTimestampSeconds();
    this.props.timeRangeChanged(now - data, now)
  }


  myAwesomeMenu() {
    let cl = this.onClick.bind(this);
    return <ContextMenu id='timeselectormenu' theme={"dark"}>
      <Item onClick={cl} data={60*15}>last 15 minutes</Item>
      <Item onClick={cl} data={60*60}>last 1 hour</Item>
      <Item onClick={cl} data={60*60*4}>last 4 hours</Item>
      <Item onClick={cl} data={60*60*12}>last 12 hours</Item>
      <Item onClick={cl} data={60*60*24}>last 24 hour</Item>
      <Item onClick={cl} data={60*60*24*7}>last 7 days</Item>
      <Item onClick={cl} data={60*60*24*30}>last 30 days</Item>
      <Item onClick={cl} data={60*60*24*90}>last 90 days</Item>
    </ContextMenu>
  }

  makeTimeseries() {
    let timeseries = this.props.logsPerTimevalue.sort((a, b) => a.TimestampPerUnit - b.TimestampPerUnit);

    let start = Math.trunc(this.props.startTime / this.props.interval);
    let finish = this.props.finishTime / this.props.interval;

    let points = [];
    let timeseriesI = 0;
    for(let x = start; x <= finish; x++) {
      if(timeseriesI < timeseries.length && timeseries[timeseriesI].TimestampPerUnit === x){
        points.push([this.props.interval + "s-" + x, timeseries[timeseriesI].Count]);
        timeseriesI++;
      }else {
        points.push([this.props.interval + "s-" + x, 0]);
      }
    }
    return points;
  }
}

export class TimeSelector extends Component {
  constructor(props) {
    super(props);
    this.state  = {
      tracker : null,
      selection : null,
    };
  }

  handleTrackerChanged = (t) => {
    this.setState({
      tracker : t,
    });
  };


  render() {
    const series = this.props.trafficSeries;
    let range = series.timerange();
    if(typeof range === 'undefined'){
      range = new TimeRange(new Date(50000), new Date(100000), )
    }

    const max = series.max();
    const tracker = this.state.tracker ? `${this.state.tracker}` : "";

    return (
      <div>
        <span className={"timeselector-trackertext"}>{tracker}</span>
        <div className="timeselector-graph">

          <Resizable>
            <ChartContainer
              timeRange={range}
              trackerPosition={this.state.tracker}
              onTrackerChanged={this.handleTrackerChanged}
              enableDragZoom={true}
              onTimeRangeChanged={
                (range) => this.props.timeRangeChanged(
                  Math.trunc(range.begin().getTime()/1000),
                  Math.trunc(range.end().getTime()/1000)
                )}
              onMouseMove={this.handleMouseMove}
              maxTime={range.end()}
              minTime={range.begin()}
              hideTimeAxis={true}
            >
              <ChartRow height="250" debug={false}>
                <YAxis
                  id="rain"
                  visible={false}
                  min={0}
                  max={max}
                  absolute={true}
                  type="linear"
                />
                <Charts>
                  <BarChart
                    axis="rain"
                    series={series}
                    style={upDownStyle}
                    colums={["value"]}
                    spacing={0}
                  />
                </Charts>

              </ChartRow>
            </ChartContainer>
          </Resizable>
        </div>
      </div>

    );
  }
}

function unixTimestampSeconds() {
  return Math.round((new Date()).getTime() / 1000)
}
