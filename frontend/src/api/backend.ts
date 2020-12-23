export class BackendAPI {

  async getRaidHealth(code: string, fight: number) {
    fetch(`/${code}/${fight}/raid_health`).then(r => r.json())
  }
}
