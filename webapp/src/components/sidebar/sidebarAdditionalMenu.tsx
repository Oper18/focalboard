// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {FormattedMessage} from 'react-intl'
import {useHistory} from 'react-router-dom'

import UsersIcon from '../../widgets/icons/users'
import VolunteerIcon from '../../widgets/icons/volunteer'

import './sidebarAdditionalMenu.scss'

const SidebarAdditionalMenu = () => {
    const history = useHistory()

    return (
        <div className='SidebarAdditionalMenu'>
            <div
                className='menu-entry'
                onClick={() => history.push('/volunteers')}
            >
                <VolunteerIcon/>
                <FormattedMessage
                    id='Sidebar.volunteers'
                    defaultMessage='Volunteers'
                />
            </div>
            <div
                className='menu-entry'
                onClick={() => history.push('/users')}
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
