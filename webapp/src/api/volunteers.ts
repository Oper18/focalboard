// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {IVolunteer} from '../types/volunteer'

import {client} from './client'

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

    const {data} = await client.get<IVolunteer[]>('/volunteers', {params})
    return data
}
