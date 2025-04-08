// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react'

import {useAppSelector} from '../../store/hooks'
import {getMe} from '../../store/users'

interface Props {
    children: React.ReactNode
}

const AdminPermissionGate = (props: Props): React.ReactElement|null => {
    const me = useAppSelector(getMe)

    // Check if user has admin role
    const isAdmin = me?.permissions?.includes('manage_system')

    if (isAdmin) {
        return <>{props.children}</>
    }
    return null
}

export default AdminPermissionGate
