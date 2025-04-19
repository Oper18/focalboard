// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
export interface IVolunteer {
    id: number
    name: string
    contact: string
    description?: string
    created_at: number
    updated_at: number
    files?: string[]
}
