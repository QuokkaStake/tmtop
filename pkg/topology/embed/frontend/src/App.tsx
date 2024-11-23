import { useState, useCallback, useEffect } from 'react'
import { fetchTopology, ComputeTopologyParams, fetchPeers, RPC } from '@/lib/api'
import { Graphviz } from 'graphviz-react'
import './App.css'

function App() {
    const [params, setParams] = useState({})
    const [dot, setDot] = useState('graph{}')
    const [peers, setPeers] = useState<RPC[]>([])
    const [selectedPeers, setSelectedPeers] = useState<{ [url: string]: boolean }>({})

    useEffect(() => {
        (async function() {
            let peers = await fetchPeers()
            setPeers(peers)
        })()
    }, [])

    const handleClick = useCallback(async () => {
        const dot = await fetchTopology({
            ...params,
            includeNodes: Object.keys(selectedPeers).filter(url => selectedPeers[url]),
        })
        setDot(dot)
    }, [params, selectedPeers, setDot, fetchTopology])

    console.log('dot', dot)

    return (
        <>
            <div style={{ display: 'flex', width: '100%' }}>
                <div style={{ 'max-width': '300px', 'overflow-x': 'scroll' }}>
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
                <div>
                    <Graphviz dot={dot} options={{ zoom: true, width: 1024, height: 1024 }} />

                    <pre><code>{dot}</code></pre>
                </div>
            </div>
        </>
    )
}

export default App
