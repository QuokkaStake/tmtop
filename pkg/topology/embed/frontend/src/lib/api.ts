import queryString from 'query-string'

export type ComputeTopologyParams = {
    includeNodes?: string[]
    crawlDistance?: number
    minBytesSec?: number
}

export async function fetchTopologyDOT(params: ComputeTopologyParams): Promise<string> {
    params.crawlDistance = params.crawlDistance || 1

    const queryParams = queryString.stringify({
        // currentHomeNode: params.currentHomeNode,
        includeNodes: params.includeNodes,
        crawlDistance: params.crawlDistance,
        minBytesSec: params.minBytesSec,
        // highlightNodes: params.highlightNodes
        format: 'dot',
    })
    console.log('params', params)
    console.log('queryParams', queryParams)

    const response = await fetch(`http://localhost:8080/topology?${queryParams}`, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
        },
    })

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
    }

    return await response.text()
}

export type JSONResponse = {
    nodes?: RPC[]
    conns?: Conn[]
}

export async function fetchTopologyJSON(params: ComputeTopologyParams): Promise<JSONResponse> {
    params.crawlDistance = params.crawlDistance || 1

    const queryParams = queryString.stringify({
        includeNodes: params.includeNodes,
        crawlDistance: params.crawlDistance,
        minBytesSec: params.minBytesSec,
    })

    const response = await fetch(`http://localhost:8080/topology?${queryParams}`, {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
        },
    })

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
    }
    return await response.json()
}

export type RPC = {
    id: string
    ip: string
    url: string
    moniker: string
    validatorAddress: string
    validatorMoniker: string
}

export type Conn = {
    from: string
    to: string
    connectionStatus: {
        send_monitor: {
            avg_rate: number
        }
        recv_monitor: {
            avg_rate: number
        }
    }
}

export async function fetchPeers(): Promise<RPC[]> {
    const response = await fetch('http://localhost:8080/peers', {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
        },
    })

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
    }

    return await response.json()
}
