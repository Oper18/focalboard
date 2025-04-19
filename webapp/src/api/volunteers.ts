// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {IVolunteer} from '../types/volunteer'
import octoClient from '../octoClient'

export const getVolunteers = async (contact?: string, limit?: number, offset?: number): Promise<IVolunteer[]> => {
    const params = new URLSearchParams()
    if (contact) {
        params.append('contact', contact)
    }
    if (limit) {
        params.append('limit', limit.toString())
    }
    if (offset) {
        params.append('offset', offset.toString())
    }

    const response = await octoClient.getJson<IVolunteer[]>(await fetch(`${octoClient.getBaseURL()}/volunteers?${params.toString()}`), [])
    return response
}
