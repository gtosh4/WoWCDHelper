import axios from 'axios'

export class BackendAPI{
  constructor(){
    this.axios = axios.create({
      baseURL: '/thcda/v1',
      timeout: 10000,
    })
  }

  async getRaidHealth(code, fight) {
    const response = await this.axios.get(`/${code}/${fight}/raid_health`, {timeout: 60000})
    return response.data
  }
}
