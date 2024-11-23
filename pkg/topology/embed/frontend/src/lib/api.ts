import queryString from 'query-string'

export type ComputeTopologyParams = {
    // currentHomeNode: string
    includeNodes?: string[]
    crawlDistance?: number
    // highlightNodes?: string[]
}

export async function fetchTopology(params: ComputeTopologyParams): Promise<string> {
    params.crawlDistance = params.crawlDistance || 1

    const queryParams = queryString.stringify({
        // currentHomeNode: params.currentHomeNode,
        includeNodes: params.includeNodes,
        crawlDistance: params.crawlDistance,
        // highlightNodes: params.highlightNodes
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

export type RPC = {
    ip: string
    url: string
    moniker: string
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
