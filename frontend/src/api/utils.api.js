import api from './config'

export const utilsApi = {
  getIfscDetails: (ifsc) => {
    return api.get(`/utils/ifsc/${ifsc}`)
  },
  
  getPincodeDetails: (pincode) => {
    return fetch(`https://api.postalpincode.in/pincode/${pincode}`).then(res => res.json())
  }
}
