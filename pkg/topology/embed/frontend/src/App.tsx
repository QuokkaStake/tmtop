import { useState, useCallback, useEffect, useRef, forwardRef } from 'react'
import { fetchTopologyJSON, ComputeTopologyParams, fetchPeers, RPC } from '@/lib/api'
import { Graphviz } from 'graphviz-react'
import { GraphCanvas, GraphCanvasRef, InternalGraphPosition, layoutProvider, recommendLayout } from 'reagraph'
import './App.css'

function App() {
    const [params, setParams] = useState({})
    const [peers, setPeers] = useState<RPC[]>([])
    const [graph, setGraph] = useState({ nodes: [], edges: [] })
    const [selectedPeers, setSelectedPeers] = useState<{ [url: string]: boolean }>({})
    const nodeRef = useRef(new Map<string, InternalGraphPosition>())
    const graphRef = useRef<GraphCanvasRef | null>(null)

    useEffect(() => {
        (async function() {
            let peers = await fetchPeers()
            setPeers(peers)
        })()
    }, [])

    const handleClick = useCallback(async () => {
        const graph = await fetchTopologyJSON({
            ...params,
            includeNodes: Object.keys(selectedPeers).filter(url => selectedPeers[url]),
        })
        console.log('graph', graph)

        const nodes = graph.nodes.map(node => ({
            id: node.id,
            label: node.moniker,
        }))

        const edges = graph.conns.map(conn => ({
            source: conn.from,
            target: conn.to,
            id: `${conn.from}-${conn.to}`,
            label: ``,
        }))

        setGraph({ nodes, edges })
    }, [params, selectedPeers, setGraph, fetchTopologyJSON])

    useEffect(() => {
        if (!graph || !graph.nodes || !graph.edges) {
            return
        }

        // let { nodes, edges } = graph
        // let layout = layoutProvider({ type: recommendLayout(nodes, edges) })
        // for (let node of nodes) {
        //     if (!nodeRef.current.has(node.id)) {
        //         let pos = layout.getNodePosition(node.id)
        //         nodeRef.current.set(node.id, pos)
        //     }
        // }
    }, [graph])

    let { nodes, edges } = graph
    let layout = recommendLayout(nodes, edges)

    return (
        <>
            <div style={{ displayxxxx: 'flex', width: '100%', height: '100vh' }}>
                <div style={{ width: '300px', height: '100vh', overflowX: 'scroll', position: 'absolute', top: 0, left: 0, zIndex: 9999, backgroundColor: '#242424' }}>
                    <button onClick={handleClick}>Generate</button>
                    <ul>
                        {peers.map(peer => (
                            <li key={peer.url}>
                                <nobr>
                                    <input
                                        type="checkbox"
                                        checked={selectedPeers[peer.url]}
                                        onChange={e => setSelectedPeers({ ...selectedPeers, [peer.url]: e.target.checked })} />
                                    {peer.moniker} ({peer.url})
                                </nobr>
                            </li>
                        ))}
                    </ul>
                </div>
                <div style={{ positionxxxx: 'relative', widthxxxx: '100%' }}>
                    {/*<Graphviz dot={dot} options={{ zoom: true, width: 1024, height: 1024 }} />*/}
                    <GraphCanvas
                        ref={graphRef}
                        nodes={nodes}
                        edges={edges}
                        draggable
                        layoutOverrides={{
                            getNodePosition: (id, nodePositionArgs) => {
                                let idx = nodes.findIndex(node => node.id === id)
                                console.log('idx', idx)
                                if (idx === -1) {
                                    idx = Math.random() * 100
                                }

                                const position = {
                                    x: 25 * idx,
                                    y: idx % 2 === 0 ? 0 : 50,
                                    z: 1,
                                }

                                return nodeRef.current?.get(id) || (function() {
                                    // This next bit is quite fraught -- do not modify unless you know what you're doing
                                    nodeRef.current.set(id, { id, x: position.x, vx: position.x, y: position.y, vy: position.y, z: 1, links: [], data: null, index: idx })
                                    return position
                                })()
                            },
                        }}
                        onNodeDragged={node => {
                            nodeRef.current.set(node.id, node.position)
                        }}
                    />

                    {/*<pre><code>{dot}</code></pre>*/}
                </div>
            </div>
        </>
    )
}

export default App
