// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {FormattedMessage} from 'react-intl'
import {useHistory} from 'react-router-dom'

import UsersIcon from '../../widgets/icons/users'
import VolunteerIcon from '../../widgets/icons/volunteer'
import {useAppSelector} from '../../store/hooks'
import {getCurrentTeamId} from '../../store/teams'
import {getMe} from '../../store/users'

import './sidebarAdditionalMenu.scss'

const SidebarAdditionalMenu = () => {
    const history = useHistory()
    const me = useAppSelector(getMe)
    const teamId = useAppSelector(getCurrentTeamId)

    // Check if user has admin permissions
    const isAdmin = me?.permissions?.includes('manage_system')

    if (!isAdmin) {
        return null
    }

    return (
        <div className='SidebarAdditionalMenu'>
            <div
                className='menu-entry'
                onClick={() => history.push(`/team/${teamId}/volunteers`)}
            >
                <VolunteerIcon/>
                <FormattedMessage
                    id='Sidebar.volunteers'
                    defaultMessage='Volunteers'
                />
            </div>
            <div
                className='menu-entry'
                onClick={() => history.push(`/team/${teamId}/users`)}
            >
                <UsersIcon/>
                <FormattedMessage
                    id='Sidebar.users'
                    defaultMessage='Users'
                />
            </div>
        </div>
    )
}

export default SidebarAdditionalMenu
