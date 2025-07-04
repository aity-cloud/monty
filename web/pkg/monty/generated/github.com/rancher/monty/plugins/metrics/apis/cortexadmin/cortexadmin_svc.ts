// @generated by service-generator v0.0.1 with parameter "target=ts,import_extension=none,ts_nocheck=false"
// @generated from file  github.com/aity-cloud/monty/plugins/metrics/apis/cortexadmin/cortexadmin.proto (package cortexadmin, syntax proto3)
/* eslint-disable */

import { ConfigRequest, ConfigResponse, DeleteRuleRequest, GetRuleRequest, LabelRequest, ListRulesRequest, ListRulesResponse, LoadRuleRequest, MatcherRequest, MetricLabels, MetricMetadata, MetricMetadataRequest, QueryRangeRequest, QueryRequest, QueryResponse, SeriesInfoList, SeriesRequest, UserIDStatsList, WriteRequest, WriteResponse } from "./cortexadmin_pb";
import { axios } from "@pkg/monty/utils/axios";
import { CortexStatus } from "./status_pb";


export async function AllUserStats(): Promise<UserIDStatsList> {
  try {
    
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/all_user_stats`
    })).data;

    const response = UserIDStatsList.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-AllUserStats:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function WriteMetrics(input: WriteRequest): Promise<WriteResponse> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-WriteMetrics:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'post',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/write_metrics`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = WriteResponse.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-WriteMetrics:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function Query(input: QueryRequest): Promise<QueryResponse> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-Query:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/query`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = QueryResponse.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-Query:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function QueryRange(input: QueryRangeRequest): Promise<QueryResponse> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-QueryRange:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/query_range`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = QueryResponse.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-QueryRange:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function GetRule(input: GetRuleRequest): Promise<QueryResponse> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-GetRule:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/rules/${input.namespace}/${input.groupName}`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = QueryResponse.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-GetRule:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function GetMetricMetadata(input: MetricMetadataRequest): Promise<MetricMetadata> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-GetMetricMetadata:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/metadata`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = MetricMetadata.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-GetMetricMetadata:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function ListRules(input: ListRulesRequest): Promise<ListRulesResponse> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-ListRules:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/rules`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = ListRulesResponse.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-ListRules:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function LoadRules(input: LoadRuleRequest): Promise<void> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-LoadRules:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'post',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/rules`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = rawResponse;
    console.info('Here is the response for a request to CortexAdmin-LoadRules:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function DeleteRule(input: DeleteRuleRequest): Promise<void> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-DeleteRule:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'delete',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/rules/${input.groupName}`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = rawResponse;
    console.info('Here is the response for a request to CortexAdmin-DeleteRule:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function FlushBlocks(): Promise<void> {
  try {
    
    const rawResponse = (await axios.request({
      method: 'post',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/flush_blocks`
    })).data;

    const response = rawResponse;
    console.info('Here is the response for a request to CortexAdmin-FlushBlocks:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function GetSeriesMetrics(input: SeriesRequest): Promise<SeriesInfoList> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-GetSeriesMetrics:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/series/metadata`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = SeriesInfoList.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-GetSeriesMetrics:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function GetMetricLabelSets(input: LabelRequest): Promise<MetricLabels> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-GetMetricLabelSets:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/series/labels`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = MetricLabels.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-GetMetricLabelSets:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function GetCortexStatus(): Promise<CortexStatus> {
  try {
    
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/status`
    })).data;

    const response = CortexStatus.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-GetCortexStatus:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function GetCortexConfig(input: ConfigRequest): Promise<ConfigResponse> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-GetCortexConfig:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/config`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = ConfigResponse.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-GetCortexConfig:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}


export async function ExtractRawSeries(input: MatcherRequest): Promise<QueryResponse> {
  try {
    
    if (input) {
      console.info('Here is the input for a request to CortexAdmin-ExtractRawSeries:', input);
    }
  
    const rawResponse = (await axios.request({
      method: 'get',
      responseType: 'arraybuffer',
      headers: {
        'Content-Type': 'application/octet-stream',
        'Accept': 'application/octet-stream',
      },
      url: `/monty-api/CortexAdmin/series/raw`,
    data: input?.toBinary() as ArrayBuffer
    })).data;

    const response = QueryResponse.fromBinary(new Uint8Array(rawResponse));
    console.info('Here is the response for a request to CortexAdmin-ExtractRawSeries:', response);
    return response;
  } catch (ex: any) {
    if (ex?.response?.data) {
      const s = String.fromCharCode.apply(null, Array.from(new Uint8Array(ex?.response?.data)));
      console.error(s);
    }
    throw ex;
  }
}

