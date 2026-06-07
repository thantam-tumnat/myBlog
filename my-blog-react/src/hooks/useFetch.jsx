import axios from 'axios';
import {useEffect, useState} from 'react';

// export default function useFetch(url) {
//     const [data,setData] =useState(null)
//     const [error,setError] =useState(null)
//     const [loading,setLoading] =useState(true)
    
//     useEffect(() => {
//         const fetchData = async ()=>{
//             setLoading(true)
//             try {
//                 const res = await fetch(url)
//                 const json = await res.json()
//                 setData(json)
//                 setLoading(false)
//             } catch (error) {
//                 setError(error)
//                 setLoading(false)
//             }
//         }
//         fetchData()
      
//     }, [url])
    
//   return {loading, error, data}
// }

export default function useFetch(url) {
    const [data,setData] =useState(null)
    const [error,setError] =useState(null)
    const [loading,setLoading] =useState(true)
    
    useEffect(() => {
        const fetchData = async ()=>{
            setLoading(true)
            try {
                // const res = await fetch(url)
                const res = await axios.get(url)
                // const json = await res.json()
                setData(res.data)
                setLoading(false)
            } catch (error) {
                setError(error)
                setLoading(false)
            }
        }
        fetchData()
      
    }, [url])
    
  return {loading, error, data}
}

